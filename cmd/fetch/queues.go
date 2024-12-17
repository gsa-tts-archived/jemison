package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GSA-TTS/jemison/config"
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
var FetchPool *pgxpool.Pool
var FetchClient *river.Client[pgx.Tx]
var FetchQueues map[string]river.QueueConfig

type FetchWorker struct {
	river.WorkerDefaults[common.FetchArgs]
}

func InitializeQueues() {
	//var fetchClient *river.Client[pgx.Tx]
	queueing.InitializeRiverQueues()

	_, fP, workers := common.CommonQueueInit()
	FetchPool = fP

	river.AddWorker(workers, &FetchWorker{})

	// Grab the number of workers from the config.
	fetchService, err := env.Env.GetUserService(ThisServiceName)
	if err != nil {
		zap.L().Error("could not fetch service config")
		log.Println(err)
		os.Exit(1)
	}
	workerCount := fetchService.GetParamInt64("workers")

	MainQueue := make(map[string]river.QueueConfig)
	MainQueue[ThisServiceName] = river.QueueConfig{MaxWorkers: int(workerCount)}

	FetchQueues = make(map[string]river.QueueConfig)
	for _, host := range config.GetListOfHosts(env.Env.AllowedHosts) {
		asciiHost := stripHostToAscii(host)
		asciiQueueName := fmt.Sprintf("fetch-%s", asciiHost)
		zap.L().Info("setting up queue", zap.String("queue_name", asciiQueueName))
		FetchQueues[asciiQueueName] = river.QueueConfig{
			MaxWorkers: 30,
		}
	}

	// Start the 'main' queue
	mainFetchClient, err := river.NewClient(riverpgxv5.New(FetchPool), &river.Config{
		Queues:  MainQueue,
		Workers: workers,
	})
	if err != nil {
		zap.L().Error("could not establish main fetch client")
		log.Println(err)
		os.Exit(1)
	}

	hostsFetchClient, err := river.NewClient(riverpgxv5.New(FetchPool), &river.Config{
		Queues:            FetchQueues,
		Workers:           workers,
		FetchCooldown:     500 * time.Millisecond,
		FetchPollInterval: 1000 * time.Millisecond,
	})
	if err != nil {
		zap.L().Error("could not establish hosts fetch client")
		log.Println(err)
		os.Exit(1)
	}

	// Start the work clients
	ctx := context.Background()
	if err := mainFetchClient.Start(ctx); err != nil {
		zap.L().Error("could not launch main fetch client.", zap.String("err", err.Error()))
		os.Exit(42)
	}
	ctx = context.Background()
	if err := hostsFetchClient.Start(ctx); err != nil {
		zap.L().Error("could not launch hosts client", zap.String("err", err.Error()))
		os.Exit(42)
	}
}
