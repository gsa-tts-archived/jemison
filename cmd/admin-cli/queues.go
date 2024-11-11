package main

import (
	"os"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"go.uber.org/zap"
)

// GLOBAL TO THE APP
var insertClient *river.Client[pgx.Tx]

func InitializeQueues() {
	queueing.InitializeRiverQueues()

	// Insert-only client
	_, pool, _ := common.CommonQueueInit()
	ic, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		zap.L().Error("could not establish insert-only client")
		os.Exit(1)
	}
	insertClient = ic
}
