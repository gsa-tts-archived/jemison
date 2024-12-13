package main

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/GSA-TTS/jemison/config"
	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/sqlite"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// We get pings on domains as they go through
// When the timer fires, we queue that domain to the finalize queue.
func FinalizeTimer(in <-chan string) {
	clocks := make(map[string]time.Time)

	// https://dev.to/milktea02/misunderstanding-go-timers-and-channels-1jal
	s, _ := env.Env.GetUserService("pack")
	TIMEOUT_DURATION := time.Duration(s.GetParamInt64("packing_timeout_seconds")) * time.Second
	zap.L().Debug("finalize starting timer")

	// How often should we check? Once per minute?
	timeout := time.NewTicker(10 * time.Second)

	for {
		select {
		case host := <-in:
			// When we see a domain, reset the clock.
			// E.g. a clock is running on alice.gov.
			// Reset it to "now" so that we have more time until it times out.
			//zap.L().Debug("resetting clock", zap.String("host", host))
			clocks[host] = time.Now()

		case <-timeout.C:
			// Every <timeout> seconds, we'll see if anyone has a clock that is greater,
			// which will mean nothing has come through recently.
			last_update := time.Now()
			next_to_pack := "NONE"

			for host, clock := range clocks {
				if time.Since(clock) > TIMEOUT_DURATION && clock.Before(last_update) {
					last_update = clock
					next_to_pack = host
				}
			}

			if next_to_pack != "NONE" {
				backend := config.GetHostBackend(next_to_pack, env.Env.AllowedHosts)
				// limit <- Ping{}
				// FIXME: This only applies if it is an SQLite DB...
				sqlite_filename := sqlite.SqliteFilename(next_to_pack)

				switch backend {
				case "sqlite":
					// Only pack one DB at a time.
					zap.L().Info("locking for pack", zap.String("host", next_to_pack))
					PackSqliteDatabase(next_to_pack)
					s3 := kv.NewS3("serve")
					err := s3.FileToS3Path(sqlite_filename, sqlite_filename, util.SQLite3.String())
					if err != nil {
						zap.L().Fatal("pack could not send file to s3",
							zap.String("sqlite_filename", sqlite_filename),
							zap.String("err", err.Error()))
					}
					os.Remove(sqlite_filename)
					zap.L().Info("unlocked", zap.String("host", next_to_pack))
				case "postgres":
					zap.L().Info("writing search data to postgres")
					PackPostgresDatabase(next_to_pack)
				}
				// <-limit

				// Enqueue next steps
				ChQSHP <- queueing.QSHP{
					Queue:    "serve",
					Filename: sqlite_filename,
				}
				delete(clocks, next_to_pack)
				runtime.GC()
			}
		}
	}
}

func (w *PackWorker) Work(ctx context.Context, job *river.Job[common.PackArgs]) error {

	ChFinalize <- job.Args.Host

	return nil
}
