package postgres

import (
	"context"
	"log"
	"time"

	"github.com/GSA-TTS/jemison/internal/env"
	search_db "github.com/GSA-TTS/jemison/internal/postgres/search_db"
	work_db "github.com/GSA-TTS/jemison/internal/postgres/work_db"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type JemisonDB struct {
	Config          map[string]*pgxpool.Config
	Pool            map[string]*pgxpool.Pool
	WorkDBQueries   *work_db.Queries
	SearchDBQueries *search_db.Queries
}

func NewJemisonDB() *JemisonDB {

	jdb := &JemisonDB{
		Config: make(map[string]*pgxpool.Config),
		Pool:   make(map[string]*pgxpool.Pool),
	}

	for _, db_name := range []string{env.QueueDatabase, env.JemisonWorkDatabase} {
		db_string, err := env.Env.GetDatabaseUrl(db_name)
		if err != nil {
			zap.L().Fatal("could not get db URL", zap.String("db_name", db_name))
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

		jdb.Config[db_name] = cfg
		jdb.Pool[db_name] = pool

	}

	jdb.WorkDBQueries = work_db.New(jdb.Pool[env.JemisonWorkDatabase])
	jdb.SearchDBQueries = search_db.New(jdb.Pool[env.SearchDatabases[0]])

	return jdb
}

func Config(db_string string) *pgxpool.Config {
	const defaultMaxConns = int32(100)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(db_string)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}
