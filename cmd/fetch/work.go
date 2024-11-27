package main

import (
	"context"
	"math/rand/v2"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	_ "embed"

	common "github.com/GSA-TTS/jemison/internal/common"
	kv "github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/GSA-TTS/jemison/internal/work_db/work_db"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// ///////////////////////////////////
// GLOBALS
var LastHitMap sync.Map
var LastBackoffMap sync.Map

var fetchCount atomic.Int64

func InfoFetchCount() {
	// Probably should be a config value.
	ticker := time.NewTicker(60 * time.Second)
	recent := []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	last := int64(0)
	ndx := 0
	for {
		// Wait for the ticker
		<-ticker.C
		cnt := fetchCount.Load()
		diff := cnt - last
		recent[ndx] = diff
		if last != 0 {
			var total int64 = 0
			for _, num := range recent {
				total += num
			}
			zap.L().Info("pages fetched",
				zap.Int64("pages", cnt),
				zap.Int64("ppm (5m avg)", total/int64(len(recent))))
		}
		ndx = (ndx + 1) % len(recent)
		last = cnt
	}
}

func randRange(min, max int) int64 {
	return int64(rand.IntN(max-min)) + int64(min)
}

func (w *FetchWorker) Work(ctx context.Context, job *river.Job[common.FetchArgs]) error {
	// This is complex. Remember, there could be 100s of workers, and we could have
	// multiple, independent `fetch` services running.
	//
	// First, we check to see if we're good to go.
	// This is a synchronized map, and it looks at the host, asking if it has come through
	// in less time than the PoliteSleep interval. If it is too soon, we are *not*
	// good to go. So, we requeue.

	if !Gateway.GoodToGo(job.Args.Host) {
		// If this host is not good to go, reschedule it for some time in the future.
		// jitter := time.Duration(randRange(-10, 10)) * time.Millisecond
		// sleepyTime := Gateway.TimeRemaining(job.Args.Host) + jitter
		// time.Sleep(sleepyTime)
		time.Sleep(time.Duration(randRange(50, 100)) * time.Millisecond)
		ChQSHP <- queueing.QSHP{
			Queue:  "fetch",
			Scheme: job.Args.Scheme,
			Host:   job.Args.Host,
			Path:   job.Args.Path,
		}
		return nil
	}

	// If we are good to go, it means that we got the lock, and the time has elapsed
	// that was needed. Why this works?
	//
	// If there are three jobs that hit at the same time for hosts A1, A2, and B1,
	// then only one will get the lock. Say A2 gets the lock. It will proceed,
	// and reset the timer for all A hosts. If A1 gets the lock, it will see that
	// it is too soon, and requeue. B1 will then get the lock, and proceed,
	// because it is not (yet) in the map, and it can clearly proceed.
	// Now, A1 will come around on the queue again, and if it is time, it
	// will proceed.

	zap.L().Debug("fetching page content", zap.String("url", host_and_path(job)))
	page_json, err := fetch_page_content(job)

	if err != nil {
		// The queueing system retries should save us here; bail if we
		// can't get the content now.
		if strings.Contains(err.Error(), common.NonIndexableContentType.String()) ||
			strings.Contains(err.Error(), common.FileTooLargeToFetch.String()) ||
			strings.Contains(err.Error(), common.FileTooSmallToProcess.String()) {
			// Return nil, because we want to consume the job, but not requeue it
			zap.L().Info("common file error", zap.String("type", err.Error()))
			return nil
		}

		// FIXME: If it is not the NonIndexable error, something else went wrong.
		// Send it back to the queue for now.
		// FIXME: we probably need a retry rate lower than 25. Log it and take a miss
		// might be a better strategy as we go live.
		return err
	}

	// Save the metadata about this page to S3
	cloudmap := kv.NewFromMap(
		ThisServiceName,
		util.ToScheme(job.Args.Scheme),
		job.Args.Host,
		job.Args.Path,
		page_json,
	)
	err = cloudmap.Save()
	// We get an error if we can't write to S3
	// This is pretty catestrophic.
	if err != nil {
		zap.L().Error("could not Save() s3json k/v",
			zap.String("key", cloudmap.Key.Render()),
		)
		return err
	}

	zap.L().Debug("stored", zap.String("key", cloudmap.Key.Render()))

	ChQSHP <- queueing.QSHP{
		Queue:  "extract",
		Scheme: job.Args.Scheme,
		Host:   job.Args.Host,
		Path:   job.Args.Path,
	}

	ChQSHP <- queueing.QSHP{
		Queue:  "walk",
		Scheme: job.Args.Scheme,
		Host:   job.Args.Host,
		Path:   job.Args.Path,
	}

	// queueing.InsertValidate("validate_fetch", common.ValidateFetchArgs{
	// 	Fetch: common.FetchArgs{
	// 		Scheme: job.Args.Scheme,
	// 		Host:   job.Args.Host,
	// 		Path:   job.Args.Path}})

	// Update the guestbook
	lastModified := time.Now()
	if v, ok := page_json["last-modified"]; ok {
		//layout := "2006-01-02 15:04:05"
		t, err := time.Parse(time.RFC1123, v)
		if err != nil {
			zap.L().Warn("could not convert last-modified")
			lastModified = time.Now()
		} else {
			lastModified = t
		}
	}
	WDB.Queries.UpdateNextFetch(work_db.FetchUpdateParams{
		Scheme:       job.Args.Scheme,
		Host:         job.Args.Host,
		Path:         job.Args.Path,
		LastModified: lastModified,
	})

	zap.L().Info("fetched", zap.String("host", job.Args.Host), zap.String("path", job.Args.Path))
	// A cute counter for the logs.
	fetchCount.Add(1)

	return nil
}
