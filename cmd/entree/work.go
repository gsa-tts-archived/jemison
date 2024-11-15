package main

import (
	"context"

	schemas_entree "github.com/GSA-TTS/jemison/cmd/entree/schemas"
	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

func guestbookCtxQueries() (context.Context, *pgx.Conn, *schemas_entree.Queries) {
	ctx := context.Background()
	db_string, err := env.Env.GetDatabaseUrl(env.DB1)
	if err != nil {
		zap.L().Fatal("could not get db URL for DB1")
	}
	conn, err := pgx.Connect(ctx, db_string)
	if err != nil {
		zap.L().Fatal("could not connect to DB1")
	}
	queries := schemas_entree.New(conn)
	return ctx, conn, queries
}

func (w *EntreeWorker) Work(ctx context.Context, job *river.Job[common.EntreeArgs]) error {
	// We're the new front door.
	// When a request comes in, we will run the algorithm described in
	// docs/design_entree.md.

	// If we have a hall pass, we'er going to send this straight over
	// to the `fetch` queue. But, we don't update any tracking, because
	// this is a special request. So, we leave the timing for the next
	// fetch alone.
	if job.Args.HallPass {
		queueing.InsertFetch(job.Args.Scheme, job.Args.Host, job.Args.Path)
		return nil
	}

	// If we don't have a hall pass, we have to make some decisions.
	// Are we past next fetch? If not, we're done.
	gbctx, conn, queries := guestbookCtxQueries()
	defer conn.Close(ctx)

	zap.L().Debug("entree host", zap.String("host", job.Args.Host))
	host_id, err := queries.GetHostId(gbctx, pgtype.Text{String: job.Args.Host, Valid: true})
	if err != nil {
		// We could get here if we try and pass a URL through that we're not tracking.
		// We should log this as an error, and log an error, but we have to eat the item
		// from the work queue. So, we return nil.
		zap.L().Info("host unknown to entree",
			zap.String("host", job.Args.Host),
			zap.String("path", job.Args.Path),
			zap.String("err", err.Error()))
		return nil
	}

	isNotYetTimeToFetch, _ := queries.ItIsNotYetTimeToFetch(ctx, pgtype.Int8{Int64: host_id, Valid: true})
	if isNotYetTimeToFetch {
		// Do nothing. Is is not  yet our time.
		return nil
	}

	// If it is time to fetch, then lets do the thing.
	queueing.InsertFetch(job.Args.Scheme, job.Args.Host, job.Args.Path)

	return nil
}
