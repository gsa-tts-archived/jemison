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

// The work client, doing the work of `ResultsApi`
var ResultsApiPool *pgxpool.Pool
var ResultsApiClient *river.Client[pgx.Tx]

type ResultsApiWorker struct {
	river.WorkerDefaults[common.ResultsApiArgs]
}

func InitializeQueues() {
	queueing.InitializeRiverQueues()

	ctx, fP, workers := common.CommonQueueInit()
	ResultsApiPool = fP

	// Essentially adds a worker "type" to the work engine.
	river.AddWorker(workers, &ResultsApiWorker{})

	// Grab the number of workers from the config.
	ResultsApiService, err := env.Env.GetUserService(ThisServiceName)
	if err != nil {
		zap.L().Error("could not ResultsApi service config")
		log.Println(err)
		os.Exit(1)
	}

	// Work client
	ResultsApiClient, err = river.NewClient(riverpgxv5.New(ResultsApiPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			ThisServiceName: {MaxWorkers: int(ResultsApiService.GetParamInt64("workers"))},
		},
		Workers: workers,
	})

	if err != nil {
		zap.L().Error("could not establish worker pool")
		log.Println(err)
		os.Exit(1)
	}

	// Start the work clients
	if err := ResultsApiClient.Start(ctx); err != nil {
		zap.L().Error("workers are not the means of production. exiting.")
		os.Exit(42)
	}
}
