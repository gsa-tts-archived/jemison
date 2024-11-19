package main

import (
	"context"
	"fmt"
	"time"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/work_db/work_db"
	"github.com/jackc/pgx/v5"
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
	host_id, err := queries.GetHostId(gbctx, job.Args.Host)
	if err != nil {
		// We could get here if we try and pass a URL through that we're not tracking.
		// We should log this as an error, and log an error, but we have to eat the item
		// from the work queue. So, we return nil.
		return -1, fmt.Errorf("host unknown to entree")
	}
	return host_id, nil
}

func updateNextFetch(host string) {
	ctx, conn, queries := guestbookCtxQueries()
	defer conn.Close(ctx)
	//host_id, _ := queries.GetHostId(ctx, host)
	queries.UpdateNextFetch(ctx, work_db.UpdateNextFetchParams{
		Column1: "weekly",
		Host:    host,
	})
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

	next_fetch, err := queries.GetHostNextFetch(gbctx, host_id)
	if err != nil {
		// If we return err, it will trigger a requeue. We don't want that.
		// We couldn't find the host in our table, so we should just skip it.
		return NotFound, nil
	}

	next_fetch_time, valid := next_fetch.TimestampValue()
	if valid == nil {
		// Are we before or after the next fetch?
		// If we are after, then go ahead and launch the crawl.
		if time.Now().After(next_fetch_time.Time) {
			zap.L().Debug("it is past time to fetch",
				zap.Int64("host_id", host_id),
				zap.String("host", host),
				zap.String("path", path),
			)
			queueing.InsertFetch(scheme, host, path)
			return DeadlinePassed, nil
		}
		return DeadlineNotYetReached, nil
	}
	return DefaultCase, nil
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
	next_fetch, err := queries.GetHostNextFetch(gbctx, host_id)
	if err != nil {
		return err
	}

	next_fetch_time, valid := next_fetch.TimestampValue()
	if valid == nil {
		// Are we before or after the next fetch?
		// If we are after, then go ahead and launch the crawl.
		if time.Now().After(next_fetch_time.Time) {
			zap.L().Info("it is past time for a full crawl",
				zap.Int64("host_id", host_id),
				zap.String("host", job.Args.Host),
				zap.String("path", job.Args.Path),
			)
			_, e := fetchOne(hallPass, host_id, job.Args.Scheme, job.Args.Host, job.Args.Path)
			if e != nil {
				zap.L().Error("error in kicking off fetchOne",
					zap.Int64("host_id", host_id),
					zap.String("host", job.Args.Host),
					zap.String("path", job.Args.Path),
				)
			}
		}

		// We must be before the next fetch. Do we have a hall pass?
		if hallPass {
			zap.L().Info("hall pass granted, enqueueing",
				zap.Int64("host_id", host_id),
				zap.String("host", job.Args.Host),
				zap.String("path", job.Args.Path),
			)
			// Set the next_fetch into the past.
			// Then, a full crawl will happen after we kick off the page we were given.
			queries.SetNextFetchToYesterday(gbctx, job.Args.Host)
			// And, kick off the page we were asked to start with.
			_, e := fetchOne(hallPass, host_id, job.Args.Scheme, job.Args.Host, job.Args.Path)
			if e != nil {
				zap.L().Error("error in kicking off fetchOne with hall pass",
					zap.Int64("host_id", host_id),
					zap.String("host", job.Args.Host),
					zap.String("path", job.Args.Path),
				)
			}
		} else {
			zap.L().Info("no hall pass, not enqueueing fetch",
				zap.Int64("host_id", host_id),
				zap.String("host", job.Args.Host),
				zap.String("path", job.Args.Path),
			)
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
		return fullCrawl(job.Args.HallPass, job)
	} else {
		_, e := fetchOne(job.Args.HallPass, host_id, job.Args.Scheme, job.Args.Host, job.Args.Path)
		return e
	}
}
