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

//go:generate stringer -type=FetchStatus
type FetchStatus int

const (
	NotPreviouslySeen FetchStatus = iota
	DeadlineNotYetReached
	DeadlinePassed
	NotFound
	HallPass
	DefaultCase
)

func fetchOne(hallPass bool, host_id int64, scheme, host, path string) (FetchStatus, error) {

	if hallPass {
		zap.L().Debug("we have a hall pass",
			zap.String("host", host),
			zap.Bool("hallpass", hallPass))
		queueing.InsertFetch(scheme, host, path)
		return HallPass, nil
	}

	// If we don't have a hall pass, we have to make some decisions.
	// Are we past next fetch? If not, we're done.
	gbctx, conn, queries := guestbookCtxQueries()
	defer conn.Close(gbctx)

	fetchStatus, _ := queries.ToFetchOrNotToFetch(gbctx, work_db.ToFetchOrNotToFetchParams{
		Host: host_id,
		Path: path,
	})

	/*
		entree               | {"level":"error","timestamp":"2024-11-16T14:11:34.973Z","caller":"entree/work.go:103","msg":"fetch default - should not get here","pid":13,"host_id":3,"host":"www.fac.gov","path":"/assets/img/updates/roadmap-08-2024.png","hallpass":false,"stacktrace":"main.fetchOne\n\t/home/jadudm/git/search/jemison/app/cmd/entree/work.go:103\nmain.(*EntreeWorker).Work\n\t/home/jadudm/git/search/jemison/app/cmd/entree/work.go:169\ngithub.com/riverqueue/river.(*wrapperWorkUnit[...]).Work\n\t/home/jadudm/go/pkg/mod/github.com/riverqueue/river@v0.13.0/work_unit_wrapper.go:30\ngithub.com/riverqueue/river.(*jobExecutor).execute.func2\n\t/home/jadudm/go/pkg/mod/github.com/riverqueue/river@v0.13.0/job_executor.go:216\ngithub.com/riverqueue/river.(*jobExecutor).execute\n\t/home/jadudm/go/pkg/mod/github.com/riverqueue/river@v0.13.0/job_executor.go:239\ngithub.com/riverqueue/river.(*jobExecutor).Execute\n\t/home/jadudm/go/pkg/mod/github.com/riverqueue/river@v0.13.0/job_executor.go:157"}
	*/
	zap.L().Debug("fetch status",
		zap.Int32("status", fetchStatus),
		zap.Int64("host_id", host_id),
		zap.String("host", host),
		zap.String("path", path),
		zap.Bool("hallpass", hallPass),
	)
	switch int32(fetchStatus) {
	case int32(NotPreviouslySeen):
		// This implies a new page that has never been crawled before,
		// becuase it has never signed the guestbook. We should fetch it.
		queueing.InsertFetch(scheme, host, path)
		return NotPreviouslySeen, nil

	case int32(DeadlineNotYetReached):
		// This is one that has been fetched recently. We don't have a hall pass,
		// so we shouldn't fetch it.
		return DeadlineNotYetReached, nil

	case int32(DeadlinePassed):
		// The last fetch is in our past. We should fetch this.
		queueing.InsertFetch(scheme, host, path)
		return DeadlinePassed, nil

	default:
		// The default is if we pass a host/path that we have never seen. If we haven't seen it, it
		// is because a full crawl (`fetch`/`walk`) decided this was non-indexable content. Note it
		// as part of debugging, but move on quietly.
		zap.L().Error("DefaultCase - should not be here")
		return DefaultCase, nil
	}
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
		_, e := fetchOne(hallPass, host_id, job.Args.Scheme, job.Args.Host, row.Path)
		if e != nil {
			// Quit if we have an error processing all of the paths.
			return e
		}
	}

	// And, kick off the page we were asked to start with.

	_, e := fetchOne(hallPass, host_id, job.Args.Scheme, job.Args.Host, job.Args.Path)
	if e != nil {
		zap.L().Error("could not kick off full crawl",
			zap.Int64("host_id", host_id),
			zap.String("host", job.Args.Host),
			zap.String("path", job.Args.Path),
		)
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

	// First, find out if we know about this domain.
	// Exit if we don't.
	host_id, err := getHostId(job)
	if err != nil {
		zap.L().Info("host unknown to entree",
			zap.String("host", job.Args.Host),
			zap.String("path", job.Args.Path),
			zap.String("err", err.Error()))
		return nil
	}

	if job.Args.FullCrawl {
		zap.L().Info("full crawl starting", zap.String("host", job.Args.Host))
		return fullCrawl(job.Args.HallPass, job)
	} else {
		_, e := fetchOne(job.Args.HallPass, host_id, job.Args.Scheme, job.Args.Host, job.Args.Path)
		return e
	}
}
