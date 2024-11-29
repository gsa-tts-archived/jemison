package main

import (
	"fmt"
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

//var fetchClient *river.Client[pgx.Tx]

// var extractPool *pgxpool.Pool

// var extractClient *river.Client[pgx.Tx]
// var walkPool *pgxpool.Pool

// var walkClient *river.Client[pgx.Tx]

type FetchWorker struct {
	river.WorkerDefaults[common.FetchArgs]
}

func InitializeQueues() {
	var fetchClient *river.Client[pgx.Tx]
	queueing.InitializeRiverQueues()

	ctx, fP, workers := common.CommonQueueInit()
	fetchPool = fP

	river.AddWorker(workers, &FetchWorker{})

	// Grab the number of workers from the config.
	fetchService, err := env.Env.GetUserService(ThisServiceName)
	if err != nil {
		zap.L().Error("could not fetch service config")
		log.Println(err)
		os.Exit(1)
	}

	// I'm going to create a round-robin pool of workers.
	// That is, if workerCount is 10, I'm going to create
	// 10x10 workers, in 10x queues. This way, I can round-robin
	// jobs to the queues. This should improve throughput significantly.
	// May have to have more configs, so this is better able to be controlled.
	workerCount := fetchService.GetParamInt64("workers")
	RoundRobinSize = workerCount

	// Start the 'main' queue
	fetchClient, err = river.NewClient(riverpgxv5.New(fetchPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			ThisServiceName: {MaxWorkers: int(workerCount)},
		},
		Workers: workers,
	})
	if err != nil {
		zap.L().Error("could not establish main worker pool")
		log.Println(err)
		os.Exit(1)
	}

	// Start the work clients
	if err := fetchClient.Start(ctx); err != nil {
		zap.L().Error("workers are not the means of production. exiting.")
		os.Exit(42)
	}

	// start the round-robin queues
	for n := range workerCount {
		queueName := fmt.Sprintf("%s-%d", ThisServiceName, n)
		fetchClient, err = river.NewClient(riverpgxv5.New(fetchPool), &river.Config{
			Queues: map[string]river.QueueConfig{
				queueName: {MaxWorkers: int(workerCount)},
			},
			// FetchCooldown: time.Duration(n*50) * time.Millisecond,
			Workers: workers,
		})
		if err != nil {
			zap.L().Error("could not establish worker pool", zap.Int64("pool number", n))
			log.Println(err)
			os.Exit(1)
		}
		// Start the work clients
		if err := fetchClient.Start(ctx); err != nil {
			zap.L().Error("workers are not the means of production. exiting.")
			os.Exit(42)
		}
	}
}
