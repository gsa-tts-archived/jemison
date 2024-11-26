package main

import (
	"os"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"go.uber.org/zap"
)

var Pool *pgxpool.Pool

type BalanceWorker struct {
	river.WorkerDefaults[common.BalanceArgs]
}

type SpaceAvailWorker struct {
	river.WorkerDefaults[common.SpaceAvailArgs]
}

func InitializeQueues() {
	queueing.InitializeRiverQueues()

	ctx, pool, workers := common.CommonQueueInit()
	Pool = pool

	// Essentially adds a worker "type" to the work engine.
	river.AddWorker(workers, &BalanceWorker{})
	river.AddWorker(workers, &SpaceAvailWorker{})

	// Grab the number of workers from the config.
	service, err := env.Env.GetUserService(ThisServiceName)
	if err != nil {
		zap.L().Error("could not fetch service config")
		os.Exit(1)
	}

	// Work client
	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			ThisServiceName: {MaxWorkers: int(service.GetParamInt64("workers"))},
		},
		Workers: workers,
	})
	if err != nil {
		zap.L().Error("could not establish worker pool")
		os.Exit(1)
	}

	// Start the work clients
	if err := client.Start(ctx); err != nil {
		zap.L().Error("workers are not the means of production. exiting.")
		os.Exit(42)
	}
}
