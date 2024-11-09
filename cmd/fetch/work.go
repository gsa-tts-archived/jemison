package main

import (
	"context"
	"net/url"
	"runtime"
	"sync"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// ///////////////////////////////////
// GLOBALS
var last_hit sync.Map
var last_backoff sync.Map

func (w *FetchWorker) Work(ctx context.Context, job *river.Job[common.FetchArgs]) error {
	// Check the cache.
	// We don't want to do anything if this is in the recently visited cache.
	zap.L().Debug("working", zap.String("url", host_and_path(job)))

	// Will aggressive GC keep us under the RAM limit?
	runtime.GC()

	cache_key := host_and_path(job)
	if _, found := recently_visited_cache.Get(cache_key); found {
		zap.L().Debug("cached", zap.String("key", cache_key))
		return nil
	} else {
		// If it is not cached, we have work to do.
		page_json, err := fetch_page_content(job)
		if err != nil {
			zap.L().Warn("could not fetch page content",
				zap.String("scheme", job.Args.Scheme),
				zap.String("host", job.Args.Host),
				zap.String("path", job.Args.Path),
			)
		}

		// FIXME
		// in the grand scheme, we may at this point want to have a queue for
		// coming back a day or two later. But, in terms of fetching... if you can't
		// get to the content... you're not going to store it. So, this
		// bails without sending it back to the queue (for now)
		if err != nil {
			u := url.URL{
				Scheme: job.Args.Scheme,
				Host:   job.Args.Host,
				Path:   job.Args.Path}
			zap.L().Info("could not fetch content; not requeueing",
				zap.String("url", u.String()))
			return nil
		}

		key := util.CreateS3Key(job.Args.Host, job.Args.Path, "json").Render()
		page_json["key"] = key

		zap.L().Debug("storing", zap.String("key", key))
		err = fetchStorage.Store(key, page_json)
		// We get an error if we can't write to S3
		if err != nil {
			zap.L().Warn("could not store k/v",
				zap.String("key", key),
			)
			return err
		}
		zap.L().Debug("stored", zap.String("key", key))

		// Update the cache
		recently_visited_cache.Set(host_and_path(job), key, 0)

		zap.L().Debug("inserting extract job")
		ctx, tx := common.CtxTx(extractPool)
		defer tx.Rollback(ctx)
		extractClient.InsertTx(ctx, tx, common.ExtractArgs{
			Key: key,
		}, &river.InsertOpts{Queue: "extract"})
		if err := tx.Commit(ctx); err != nil {
			zap.L().Panic("cannot commit insert tx",
				zap.String("key", key))
		}

		zap.L().Debug("Inserting walk job")
		ctx2, tx2 := common.CtxTx(walkPool)
		defer tx2.Rollback(ctx)
		walkClient.InsertTx(ctx2, tx2, common.WalkArgs{
			Key: key,
		}, &river.InsertOpts{Queue: "walk"})
		if err := tx2.Commit(ctx2); err != nil {
			zap.L().Panic("cannot commit insert tx",
				zap.String("key", key), zap.String("error", err.Error()))
		}
	}

	zap.L().Debug("Inserting validate job")
	ctx2, tx2 := common.CtxTx(walkPool)
	defer tx2.Rollback(ctx)
	walkClient.InsertTx(ctx2, tx2, common.ValidateFetchArgs{
		Fetch: common.FetchArgs{
			Scheme: job.Args.Scheme,
			Host:   job.Args.Host,
			Path:   job.Args.Path,
		},
	}, &river.InsertOpts{Queue: "validate_fetch"})
	if err := tx2.Commit(ctx2); err != nil {
		zap.L().Panic("cannot commit validate tx",
			zap.String("error", err.Error()))
	}

	zap.L().Info("fetched",
		zap.String("scheme", job.Args.Scheme),
		zap.String("host", job.Args.Host),
		zap.String("path", job.Args.Path))

	return nil
}
