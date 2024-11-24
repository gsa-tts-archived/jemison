package main

import (
	"context"
	"log"
	"time"

	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/work_db/work_db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var WDB *WorkDB

type WorkDB struct {
	Queries *work_db.Queries
	Config  *pgxpool.Config
	Pool    *pgxpool.Pool
	//DB      *pgx.Conn
}

func NewGuestbookDB() *WorkDB {
	db_string, err := env.Env.GetDatabaseUrl(env.JemisonWorkDatabase)
	if err != nil {
		zap.L().Fatal("could not get db URL for work-db")
	}

	cfg := Config(db_string)
	// Create database connection
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		zap.L().Fatal("could not create pool",
			zap.String("err", err.Error()))
	}
	// db, err := pgx.Connect(context.Background(), db_string)
	// if err != nil {
	// 	zap.L().Fatal("Error while creating connection to the database")
	// }
	err = pool.Ping(context.Background())
	if err != nil {
		zap.L().Error(err.Error())
	}

	// conn, err := pgx.Connect(ctx, db_string)
	// if err != nil {
	// 	zap.L().Fatal("could not connect to work-db")
	// }
	queries := work_db.New(pool)

	return &WorkDB{
		Queries: queries,
		Config:  cfg,
		//DB:      db,
		Pool: pool,
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

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		zap.L().Debug("BeforeAquire")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		zap.L().Debug("AfterRelease")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		zap.L().Debug("BeforeClose")
	}

	return dbConfig
}
