package work_db

import (
	"context"
	"time"

	"github.com/GSA-TTS/jemison/config"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type FetchUpdateParams struct {
	Scheme        string
	Host          string
	Path          string
	LastModified  time.Time
	ContentSha1   string
	ContentLength int64
	ContentType   string
}

func GetWorkDbQueryContext() (context.Context, *pgx.Conn, *Queries) {
	ctx := context.Background()
	db_string, err := env.Env.GetDatabaseUrl(env.JemisonWorkDatabase)
	if err != nil {
		zap.L().Fatal("could not get db URL for DB1")
	}
	conn, err := pgx.Connect(ctx, db_string)
	if err != nil {
		zap.L().Fatal("could not connect to DB1")
	}
	queries := New(conn)
	return ctx, conn, queries
}

func (q *Queries) UpdateNextFetch(params FetchUpdateParams) {
	// ctx, conn, queries := GetWorkDbQueryContext()
	// defer conn.Close(ctx)
	ctx := context.Background()
	host_id, err := q.GetHostId(ctx, params.Host)
	if err != nil {
		zap.L().Fatal("cannot get host id to update next fetch",
			zap.Int64("host_id", host_id),
			zap.String("host", params.Host),
			zap.String("path", params.Path))
	}

	schedule := config.GetScheduleFromHost(params.Host, env.Env.Schedule)

	zap.L().Debug("schedule for host",
		zap.String("host", params.Host),
		zap.String("schedule", schedule))

	// Prep for nulls...
	contentsha1 := pgtype.Text{
		Valid: false,
	}
	if len(params.ContentSha1) > 0 {
		contentsha1 = pgtype.Text{
			String: params.ContentSha1,
			Valid:  true,
		}
	}

	contenttype := pgtype.Text{
		Valid: false,
	}
	if len(params.ContentSha1) > 0 {
		contenttype = pgtype.Text{
			String: params.ContentType,
			Valid:  true,
		}
	}

	contentlength := pgtype.Int4{
		Valid: false,
	}
	if params.ContentLength != 0 {
		contentlength = pgtype.Int4{
			Valid: true,
			Int32: int32(params.ContentLength),
		}
	}

	next_fetch := config.HostToPgTimestamp(params.Host, env.Env.Schedule, time.Now())
	q.UpdateGuestbookFetch(ctx, UpdateGuestbookFetchParams{
		Scheme: params.Scheme,
		Host:   host_id,
		Path:   params.Path,
		LastUpdated: pgtype.Timestamp{
			Time:             params.LastModified,
			InfinityModifier: 0,
			Valid:            true,
		},
		ContentSha1:   contentsha1,
		ContentLength: contentlength,
		ContentType:   contenttype,
		NextFetch:     next_fetch,
	})
}

/*
	arg.Scheme,
	arg.Host,
	arg.Path,
	arg.ContentSha1,
	arg.ContentLength,
	arg.ContentType,
	arg.LastUpdated,
	arg.NextFetch,
*/
