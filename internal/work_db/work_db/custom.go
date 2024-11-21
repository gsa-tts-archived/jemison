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
	Scheme      string
	Host        string
	Path        string
	LastUpdated time.Time
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

func UpdateNextFetch(params FetchUpdateParams) {
	ctx, conn, queries := GetWorkDbQueryContext()
	defer conn.Close(ctx)
	host_id, err := queries.GetHostId(ctx, params.Host)
	if err != nil {
		zap.L().Fatal("cannot get host id to update next fetch",
			zap.Int64("host_id", host_id),
			zap.String("host", params.Host),
			zap.String("path", params.Path))
	}
	//host_id, _ := queries.GetHostId(ctx, host)
	schedule := config.GetScheduleFromHost(params.Host)

	zap.L().Debug("schedule for host",
		zap.String("host", params.Host),
		zap.String("schedule", schedule))

	next_fetch := config.HostToPgTimestamp(params.Host, time.Now())
	queries.UpdateGuestbookFetch(ctx, UpdateGuestbookFetchParams{
		Scheme: params.Scheme,
		Host:   host_id,
		Path:   params.Path,
		LastUpdated: pgtype.Timestamp{
			Time:             params.LastUpdated,
			InfinityModifier: 0,
			Valid:            true,
		},
		NextFetch: next_fetch,
	})
}
