package main

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/sqlite"
	"github.com/GSA-TTS/jemison/internal/sqlite/schemas"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/riverqueue/river"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var databases sync.Map

var ch_finalize chan *sqlite.PackTable

// We get pings on domains as they go through
// When the timer fires, we queue that domain to the finalize queue.
func FinalizeTimer(in <-chan *sqlite.PackTable) {
	clocks := make(map[string]time.Time)
	tables := make(map[string]*sqlite.PackTable)

	// https://dev.to/milktea02/misunderstanding-go-timers-and-channels-1jal
	s, _ := env.Env.GetUserService("pack")
	TIMEOUT_DURATION := time.Duration(s.GetParamInt64("packing_timeout_seconds")) * time.Second
	zap.L().Debug("finalize starting timer")
	timeout := time.NewTimer(TIMEOUT_DURATION)

	for {
		select {
		case pt := <-in:
			// When we get a domain, we should indicate that we
			// saw it just now.
			clocks[pt.Filename] = time.Now()
			tables[pt.Filename] = pt

		case <-timeout.C:
			// Every <timeout> seconds, we'll see if anyone has a clock that is greater,
			// which will mean nothing has come through recently.
			//zap.L().Debug("finalize timeout")
			for sqlite_filename, clock := range clocks {
				if time.Since(clock) > TIMEOUT_DURATION {
					zap.L().Info("packing to sqlite",
						zap.String("sqlite_filename", sqlite_filename))

					tables[sqlite_filename].PrepForNetwork()

					s3 := kv.NewS3("serve")
					err := s3.FileToS3Path(sqlite_filename, sqlite_filename, util.SQLite3.String())
					if err != nil {
						log.Println("PACK could not store to file", sqlite_filename)
						log.Fatal(err)
					}

					// Enqueue serve
					zap.L().Debug("inserting serve job")
					ctx, tx := common.CtxTx(servePool)
					serveClient.InsertTx(ctx, tx, common.ServeArgs{
						Filename: sqlite_filename,
					}, &river.InsertOpts{Queue: "serve"})
					if err := tx.Commit(ctx); err != nil {
						tx.Rollback(ctx)
						zap.L().Panic("cannot commit insert tx",
							zap.String("filename", sqlite_filename))
					}

					delete(clocks, sqlite_filename)
					delete(tables, sqlite_filename)

				}
			}
		}
		//zap.L().Debug("finalize reset")
		timeout.Reset(TIMEOUT_DURATION)
	}
}

func (w *PackWorker) Work(ctx context.Context, job *river.Job[common.PackArgs]) error {
	zap.L().Debug("packing")

	s3json := kv.NewEmptyS3JSON("extract",
		util.ToScheme(job.Args.Scheme),
		job.Args.Host,
		job.Args.Path)
	s3json.Load()

	host := s3json.GetString("host")

	JSON := s3json.GetJSON()
	pt, err := sqlite.CreatePackTable(sqlite.SqliteFilename(host), JSON)
	if err != nil {
		log.Println("Could not create pack table for", host)
		log.Fatal(err)
	}

	_, err = pt.Queries.CreateSiteEntry(pt.Context, schemas.CreateSiteEntryParams{
		Host: gjson.GetBytes(JSON, "host").String(),
		Path: gjson.GetBytes(JSON, "path").String(),
		Text: gjson.GetBytes(JSON, "content").String(),
	})
	if err != nil {
		log.Println("Insert into site entry table failed")
		log.Fatal(err)
	}

	zap.L().Debug("packed entry",
		zap.String("database", host),
		zap.String("path", gjson.GetBytes(JSON, "path").String()),
		zap.Int("length", len(gjson.GetBytes(JSON, "content").String())))

	pt.DB.Close()

	ch_finalize <- pt

	// Agressively keep memory clear.
	// GC after packing every message.
	runtime.GC()
	return nil
}
