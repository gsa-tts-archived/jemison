package search_db

import (
	"context"
	"log"
	"time"

	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SearchDB struct {
	Index   string
	Queries *Queries
	Config  *pgxpool.Config
	Pool    *pgxpool.Pool
}

func NewSearchDB(searchDB string) *SearchDB {
	db_string, err := env.Env.GetDatabaseUrl(searchDB)
	if err != nil {
		zap.L().Fatal("could not get db URL for ", zap.String("searchDB", searchDB))
	}

	cfg := Config(db_string)
	// Create database connection
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		zap.L().Fatal("could not create pool",
			zap.String("err", err.Error()))
	}

	err = pool.Ping(context.Background())
	if err != nil {
		zap.L().Error(err.Error())
	}

	queries := New(pool)

	return &SearchDB{
		Index:   searchDB,
		Queries: queries,
		Config:  cfg,
		Pool:    pool,
	}
}

func Config(db_string string) *pgxpool.Config {
	const defaultMaxConns = int32(100)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Your own Database URL
	var DATABASE_URL string = db_string

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	// dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
	// 	zap.L().Debug("BeforeAquire")
	// 	return true
	// }

	// dbConfig.AfterRelease = func(c *pgx.Conn) bool {
	// 	//zap.L().Debug("AfterRelease")
	// 	return true
	// }

	// dbConfig.BeforeClose = func(c *pgx.Conn) {
	// 	//zap.L().Debug("BeforeClose")
	// }

	return dbConfig
}
