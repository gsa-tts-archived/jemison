package main

import (
	"log"
	"os"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"go.uber.org/zap"
)

// The work client, doing the work of `Collect`.
var collectPool *pgxpool.Pool

var collectClient *river.Client[pgx.Tx]

type CollectWorker struct {
	river.WorkerDefaults[common.CollectArgs]
}

func InitializeQueues() {
	queueing.InitializeRiverQueues()

	ctx, fP, workers := common.CommonQueueInit()
	collectPool = fP

	// Essentially adds a worker "type" to the work engine.
	river.AddWorker(workers, &CollectWorker{})

	// Grab the number of workers from the config.
	collectService, err := env.Env.GetUserService(ThisServiceName)
	if err != nil {
		zap.L().Error("could not ResultsApi service config")
		log.Println(err)
		os.Exit(1)
	}

	// Work client
	collectClient, err = river.NewClient(riverpgxv5.New(collectPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			ThisServiceName: {MaxWorkers: int(collectService.GetParamInt64("workers"))},
		},
		Workers: workers,
	})

	if err != nil {
		zap.L().Error("could not establish worker pool")
		log.Println(err)
		os.Exit(1)
	}

	// Start the work clients
	if err := collectClient.Start(ctx); err != nil {
		zap.L().Error("workers are not the means of production. exiting.")
		os.Exit(1)
	}
}
