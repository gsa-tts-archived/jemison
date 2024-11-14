package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
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

func getSpaceAvailable() uint64 {
	var stat syscall.Statfs_t

	err := syscall.Statfs("/", &stat)
	if err != nil {
		zap.L().Fatal("could not get disk space for packing")
	}

	// Available space in bytes
	availableBytes := stat.Bavail * uint64(stat.Bsize)
	return availableBytes
}

func optionalSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	} else {
		return path + "/"
	}
}

func PackTheDatabase(host string) {

	// Look into the "extract" bucket and get a
	// list of all the objects there.
	s3 := kv.NewS3("extract")

	list_of_objects, err := s3.List(optionalSlash(host))
	if err != nil {
		zap.L().Fatal("could not list objects",
			zap.String("bucket", "extract"),
			zap.String("host", host))
	}

	// Calculate size
	var size int64 = 0
	for _, s3obj := range list_of_objects {
		zap.L().Debug("size", zap.String("key", s3obj.Key), zap.Int64("size", s3obj.Size))
		size += s3obj.Size
	}
	available := getSpaceAvailable()
	if (float64(size) * 1.1) > float64(available) {
		// If we have more to pack than is available on the machine
		// 1) we're in trouble
		// 2) we want to send this back to the queue.
		//
		// FIXME: Add a client to queue this back for ourselves.
		return
	}

	pt, err := sqlite.CreatePackTable(sqlite.SqliteFilename(host))
	if err != nil {
		log.Println("Could not create pack table for", host)
		log.Fatal(err)
	}

	for _, s3obj := range list_of_objects {
		s3json, err := s3.S3PathToS3JSON(s3obj.Key)
		if err != nil {
			zap.L().Error("could not fetch object for packing",
				zap.String("key", s3json.Key.Render()))
		}
		JSON := s3json.GetJSON()
		_, err = pt.Queries.CreateSiteEntry(pt.Context, schemas.CreateSiteEntryParams{
			Host: gjson.GetBytes(JSON, "host").String(),
			Path: gjson.GetBytes(JSON, "path").String(),
			Text: gjson.GetBytes(JSON, "content").String(),
		})
		if err != nil {
			log.Println("Insert into site entry table failed")
			log.Fatal(err)
		}
	}

	// Cleanup time.
	pt.DB.Close()
	pt.PrepForNetwork()
}

// We get pings on domains as they go through
// When the timer fires, we queue that domain to the finalize queue.
func FinalizeTimer(in <-chan string) {
	clocks := make(map[string]time.Time)

	// https://dev.to/milktea02/misunderstanding-go-timers-and-channels-1jal
	s, _ := env.Env.GetUserService("pack")
	TIMEOUT_DURATION := time.Duration(s.GetParamInt64("packing_timeout_seconds")) * time.Second
	zap.L().Debug("finalize starting timer")
	timeout := time.NewTimer(TIMEOUT_DURATION)

	for {
		select {
		case host := <-in:
			// When we see a domain, reset the clock.
			// E.g. a clock is running on alice.gov.
			// Reset it to "now" so that we have more time until it times out.
			zap.L().Debug("resetting clock", zap.String("host", host))
			clocks[host] = time.Now()

		case <-timeout.C:
			// Every <timeout> seconds, we'll see if anyone has a clock that is greater,
			// which will mean nothing has come through recently.
			for host, clock := range clocks {
				if time.Since(clock) > TIMEOUT_DURATION {

					PackTheDatabase(host)
					sqlite_filename := sqlite.SqliteFilename(host)
					s3 := kv.NewS3("serve")
					err := s3.FileToS3Path(sqlite_filename, sqlite_filename, util.SQLite3.String())
					if err != nil {
						log.Println("PACK could not store to file", sqlite_filename)
						log.Fatal(err)
					}
					os.Remove(sqlite_filename)

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

					delete(clocks, host)
					runtime.GC()
				}
			}
		}
		//zap.L().Debug("finalize reset")
		timeout.Reset(TIMEOUT_DURATION)
	}
}

func (w *PackWorker) Work(ctx context.Context, job *river.Job[common.PackArgs]) error {

	ChFinalize <- job.Args.Host

	return nil
}
