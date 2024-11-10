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

// GLOBAL TO THE APP
// One pool of connections for River.

// The work client, doing the work of `fetch`
var fetchPool *pgxpool.Pool
var fetchClient *river.Client[pgx.Tx]
var extractPool *pgxpool.Pool
var extractClient *river.Client[pgx.Tx]
var walkPool *pgxpool.Pool
var walkClient *river.Client[pgx.Tx]

type FetchWorker struct {
	river.WorkerDefaults[common.FetchArgs]
}

func InitializeQueues() {
	queueing.InitializeRiverQueues()

	ctx, fP, workers := common.CommonQueueInit()
	_, eP, _ := common.CommonQueueInit()
	_, wP, _ := common.CommonQueueInit()
	fetchPool = fP
	extractPool = eP
	walkPool = wP

	// Essentially adds a worker "type" to the work engine.
	river.AddWorker(workers, &FetchWorker{})

	// Grab the number of workers from the config.
	fetchService, err := env.Env.GetUserService(ThisServiceName)
	if err != nil {
		zap.L().Error("could not fetch service config")
		log.Println(err)
		os.Exit(1)
	}

	// Work client
	fetchClient, err = river.NewClient(riverpgxv5.New(fetchPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			ThisServiceName: {MaxWorkers: int(fetchService.GetParamInt64("workers"))},
		},
		Workers: workers,
	})

	if err != nil {
		zap.L().Error("could not establish worker pool")
		log.Println(err)
		os.Exit(1)
	}

	// Insert-only client to `extract`
	extractClient, err = river.NewClient(riverpgxv5.New(extractPool), &river.Config{})
	if err != nil {
		zap.L().Error("could not establish insert-only client")
		os.Exit(1)
	}
	walkClient, err = river.NewClient(riverpgxv5.New(walkPool), &river.Config{})
	if err != nil {
		zap.L().Error("could not establish insert-only client")
		os.Exit(1)
	}

	// Start the work clients
	if err := fetchClient.Start(ctx); err != nil {
		zap.L().Error("workers are not the means of production. exiting.")
		os.Exit(42)
	}
}
