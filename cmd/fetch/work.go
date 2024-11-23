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
var last_hit sync.Map
var last_backoff sync.Map

var fetchCount atomic.Int64

func InfoFetchCount() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		<-ticker.C
		cnt := fetchCount.Load()
		zap.L().Info("pages fetched",
			zap.Int64("pages", cnt))
	}
}

func (w *FetchWorker) Work(ctx context.Context, job *river.Job[common.FetchArgs]) error {
	// Check the cache.
	// We don't want to do anything if this is in the recently visited cache.

	// Use the queue instead of sleeping. It is keeping fetchers from working.
	last_hit_time, ok := last_hit.Load(job.Args.Host)
	// If we're in the map, and we're within 2s, we should keep checking after a backoff
	polite_duration := time.Duration(polite_sleep) * time.Second

	if ok && (time.Since(last_hit_time.(time.Time)) < polite_duration) {
		ch_qshp <- queueing.QSHP{
			Queue:  "fetch",
			Scheme: job.Args.Scheme,
			Host:   job.Args.Host,
			Path:   job.Args.Path,
		}
		// queueing.InsertFetch(
		// 	job.Args.Scheme,
		// 	job.Args.Host,
		// 	job.Args.Path,
		// )

		// Put a *less tiny* backoff. We'll see how the queue does.
		time.Sleep(time.Duration(rand.IntN(50)+1) * time.Millisecond)
		return nil
	}

	zap.L().Debug("fetching page content", zap.String("url", host_and_path(job)))
	page_json, err := fetch_page_content(job)
	if err != nil {
		// The queueing system retries should save us here; bail if we
		// can't get the content now.
		if strings.Contains(err.Error(), common.NonIndexableContentType.String()) {
			zap.L().Info("could not fetch page content",
				zap.String("scheme", job.Args.Scheme),
				zap.String("host", job.Args.Host),
				zap.String("path", job.Args.Path),
			)
			// Return nil, because we want to consume the job, but not requeue it
			return nil
		}
		// Otherwise... requeue?
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
	if err != nil {
		zap.L().Warn("could not Save() s3json k/v",
			zap.String("key", cloudmap.Key.Render()),
		)
		return err
	}
	zap.L().Debug("stored", zap.String("key", cloudmap.Key.Render()))

	// queueing.InsertExtract(
	// 	job.Args.Scheme,
	// 	job.Args.Host,
	// 	job.Args.Path,
	// )
	ch_qshp <- queueing.QSHP{
		Queue:  "extract",
		Scheme: job.Args.Scheme,
		Host:   job.Args.Host,
		Path:   job.Args.Path,
	}

	// queueing.InsertWalk(
	// 	job.Args.Scheme,
	// 	job.Args.Host,
	// 	job.Args.Path,
	// )
	ch_qshp <- queueing.QSHP{
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
	last_updated := time.Now()
	if v, ok := page_json["last-updated"]; ok {
		layout := "2006-01-02 15:04:05"
		t, err := time.Parse(layout, v)
		if err != nil {
			zap.L().Warn("could not convert time for last-updated")
			last_updated = time.Now()
		} else {
			last_updated = t
		}
	}

	work_db.UpdateNextFetch(work_db.FetchUpdateParams{
		Scheme:      job.Args.Scheme,
		Host:        job.Args.Host,
		Path:        job.Args.Path,
		LastUpdated: last_updated,
	})

	fetchCount.Add(1)

	zap.L().Debug("fetched",
		zap.String("scheme", job.Args.Scheme),
		zap.String("host", job.Args.Host),
		zap.String("path", job.Args.Path),
	)

	return nil
}
