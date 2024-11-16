package main

import (
	"context"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	_ "embed"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	kv "github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/GSA-TTS/jemison/internal/work_db/work_db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

func guestbookCtxQueries() (context.Context, *pgx.Conn, *work_db.Queries) {
	ctx := context.Background()
	db_string, err := env.Env.GetDatabaseUrl(env.JemisonWorkDatabase)
	if err != nil {
		zap.L().Fatal("could not get db URL for DB1")
	}
	conn, err := pgx.Connect(ctx, db_string)
	if err != nil {
		zap.L().Fatal("could not connect to DB1")
	}
	queries := work_db.New(conn)
	return ctx, conn, queries
}

func updateGuestbook(scheme, host, path string) {
	ctx, conn, queries := guestbookCtxQueries()
	defer conn.Close(ctx)
	host_id, _ := queries.GetHostId(ctx, pgtype.Text{String: host, Valid: true})
	queries.UpdateNextFetch(ctx, work_db.UpdateNextFetchParams{
		Column1: "weekly",
		Scheme:  scheme,
		Host:    host_id,
		Path:    path,
	})
}

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
	zap.L().Debug("working", zap.String("url", host_and_path(job)))

	// Will aggressive GC keep us under the RAM limit?
	runtime.GC()

	cache_key := host_and_path(job)
	if _, found := RecentlyVisitedCache.Get(cache_key); found {
		zap.L().Debug("cached", zap.String("key", cache_key))
		return nil
	} else {
		// If it is not cached, we have work to do.
		page_json, err := fetch_page_content(job)
		if err != nil {
			if !strings.Contains(err.Error(), common.NonIndexableContentType.String()) {
				zap.L().Warn("could not fetch page content",
					zap.String("scheme", job.Args.Scheme),
					zap.String("host", job.Args.Host),
					zap.String("path", job.Args.Path),
				)
			}
		}

		// The queueing system retries should save us here; bail if we
		// can't get the content now.
		if err != nil {
			u := url.URL{
				Scheme: job.Args.Scheme,
				Host:   job.Args.Host,
				Path:   job.Args.Path}
			zap.L().Debug("could not fetch content; not requeueing",
				zap.String("url", u.String()))
			return nil
		}

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

		// Update the cache
		RecentlyVisitedCache.Set(host_and_path(job), cloudmap.Key.Render(), 0)

		zap.L().Debug("inserting extract job")
		ctx, tx := common.CtxTx(extractPool)
		defer tx.Rollback(ctx)
		extractClient.InsertTx(ctx, tx, common.ExtractArgs{
			Scheme: job.Args.Scheme,
			Host:   job.Args.Host,
			Path:   job.Args.Path,
		}, &river.InsertOpts{Queue: "extract"})
		if err := tx.Commit(ctx); err != nil {
			zap.L().Panic("cannot commit insert tx",
				zap.String("key", cloudmap.Key.Render()))
		}

		zap.L().Debug("Inserting walk job")
		ctx2, tx2 := common.CtxTx(walkPool)
		defer tx2.Rollback(ctx)
		walkClient.InsertTx(ctx2, tx2, common.WalkArgs{
			Scheme: job.Args.Scheme,
			Host:   job.Args.Host,
			Path:   job.Args.Path,
		}, &river.InsertOpts{Queue: "walk"})
		if err := tx2.Commit(ctx2); err != nil {
			zap.L().Panic("cannot commit insert tx",
				zap.String("key", cloudmap.Key.Render()), zap.String("error", err.Error()))
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

	// Update the guestbook
	updateGuestbook(job.Args.Scheme, job.Args.Host, job.Args.Path)

	fetchCount.Add(1)

	zap.L().Debug("fetched",
		zap.String("scheme", job.Args.Scheme),
		zap.String("host", job.Args.Host),
		zap.String("path", job.Args.Path),
	)

	return nil
}
