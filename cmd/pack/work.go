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
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/sqlite"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/riverqueue/river"
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
		zap.L().Warn("not enough space to pack the database",
			zap.Float64("database size", float64(size)), zap.Float64("available", float64(available)))
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

		contentType := s3json.GetString("content-type")
		switch contentType {
		case "text/html":
			packHtml(pt, s3json)
		case "application/pdf":
			packPdf(pt, s3json)
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
			//zap.L().Debug("resetting clock", zap.String("host", host))
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
						zap.L().Fatal("pack could not store to file",
							zap.String("sqlite_filename", sqlite_filename),
							zap.String("err", err.Error()))
					}
					os.Remove(sqlite_filename)

					// Enqueue next steps
					ChQSHP <- queueing.QSHP{
						Queue:    "serve",
						Filename: sqlite_filename,
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
