package main

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	_ "embed"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	kv "github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/queueing"
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

	page_json, err := fetch_page_content(job)
	if err != nil {
		// The queueing system retries should save us here; bail if we
		// can't get the content now.
		if !strings.Contains(err.Error(), common.NonIndexableContentType.String()) {
			zap.L().Warn("could not fetch page content",
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

	queueing.InsertExtract(
		job.Args.Scheme,
		job.Args.Host,
		job.Args.Path,
	)

	queueing.InsertWalk(
		job.Args.Scheme,
		job.Args.Host,
		job.Args.Path,
	)

	queueing.InsertValidate("validate_fetch", common.ValidateFetchArgs{
		Fetch: common.FetchArgs{
			Scheme: job.Args.Scheme,
			Host:   job.Args.Host,
			Path:   job.Args.Path}})

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
