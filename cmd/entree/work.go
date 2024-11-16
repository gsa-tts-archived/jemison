package main

import (
	"context"
	"fmt"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
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

func getHostId(job *river.Job[common.EntreeArgs]) (int64, error) {
	gbctx, conn, queries := guestbookCtxQueries()
	defer conn.Close(gbctx)
	host_id, err := queries.GetHostId(gbctx, pgtype.Text{String: job.Args.Host, Valid: true})
	if err != nil {
		// We could get here if we try and pass a URL through that we're not tracking.
		// We should log this as an error, and log an error, but we have to eat the item
		// from the work queue. So, we return nil.
		return -1, fmt.Errorf("host unknown to entree")
	}
	return host_id, nil
}

func fetchOne(hallPass bool, host_id int64, scheme, host, path string) error {

	if hallPass {
		zap.L().Debug("we have a hall pass",
			zap.String("host", host),
			zap.Bool("hallpass", hallPass))
		queueing.InsertFetch(scheme, host, path)
		return nil
	}

	// If we don't have a hall pass, we have to make some decisions.
	// Are we past next fetch? If not, we're done.
	gbctx, conn, queries := guestbookCtxQueries()
	defer conn.Close(gbctx)

	fetchStatus, _ := queries.ToFetchOrNotToFetch(gbctx, work_db.ToFetchOrNotToFetchParams{
		Host: host_id,
		Path: path,
	})

	switch fetchStatus {
	case "FETCH_NOT_FOUND":
		// This implies a new page that has never been crawled before,
		// becuase it has never signed the guestbook. We should fetch it.
		zap.L().Debug("fetch uncertain",
			zap.String("host", host),
			zap.Bool("hallpass", hallPass),
		)
		queueing.InsertFetch(scheme, host, path)
		return nil

	case "DO_NOT_FETCH":
		// This is one that has been fetched recently. We don't have a hall pass,
		// so we shouldn't fetch it.
		return nil

	case "FETCH":
		// The last fetch is in our past. We should fetch this.
		zap.L().Debug("it is time to fetch",
			zap.String("host", host),
			zap.Bool("hallpass", hallPass),
		)
		queueing.InsertFetch(scheme, host, path)
		return nil

	default:
		zap.L().Error("found no fetch status",
			zap.String("host", host),
			zap.String("path", host),
			zap.Bool("hallpass", hallPass))
	}

	return nil
}

func fullCrawl(hallPass bool, job *river.Job[common.EntreeArgs]) error {
	// Get all of the paths in this domain
	gbctx, conn, queries := guestbookCtxQueries()
	defer conn.Close(gbctx)
	host_id, err := getHostId(job)
	if err != nil {
		zap.L().Info("host unknown to entree",
			zap.String("host", job.Args.Host),
			zap.String("path", job.Args.Path),
			zap.String("err", err.Error()))
		return nil
	}

	// if a full crawl is requested, we should look at all the existing pages.
	// potentially visit each, using our hall pass.
	rows, err := queries.GetGuestbookEntries(gbctx, host_id)
	if err != nil {
		return err
	}

	// // Given all the paths for the host, check each one.
	for _, row := range rows {
		e := fetchOne(hallPass, host_id, job.Args.Scheme, job.Args.Host, row.Path)
		if e != nil {
			// Quit if we have an error processing all of the paths.
			return e
		}
	}

	return nil
}

func (w *EntreeWorker) Work(ctx context.Context, job *river.Job[common.EntreeArgs]) error {
	// We're the new front door.
	// When a request comes in, we will run the algorithm described in
	// docs/design_entree.md.

	// Matrix
	// fullCrawl & !pass: check every timeout in the domain.
	// fullCrawl & pass: re-crawl the whole domain now.
	// !fullCrawl & !pass: check
	// !fullCrawl & pass: fetch the page now

	if job.Args.FullCrawl {
		return fullCrawl(job.Args.HallPass, job)
	} else {
		host_id, err := getHostId(job)
		if err != nil {
			zap.L().Info("host unknown to entree",
				zap.String("host", job.Args.Host),
				zap.String("path", job.Args.Path),
				zap.String("err", err.Error()))
			return nil
		}
		return fetchOne(job.Args.HallPass, host_id, job.Args.Scheme, job.Args.Host, job.Args.Path)
	}
}
