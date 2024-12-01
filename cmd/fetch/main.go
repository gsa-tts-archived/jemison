package main

import (
	"log"
	"net/http"
	"time"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/work_db/work_db"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

var PoliteSleep int64
var ThisServiceName = "fetch"

var RetryClient *http.Client
var Gateway *HostGateway

var WDB *work_db.WorkDB
var QDB *work_db.QueueDB

var ChQSHP = make(chan queueing.QSHP)

var Workers *river.Workers

func main() {
	env.InitGlobalEnv(ThisServiceName)
	InitializeQueues()

	engine := common.InitializeAPI()
	ExtendApi(engine)

	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryMax = 10
	RetryClient = retryableClient.StandardClient()

	log.Println("environment initialized")

	// Init a cache for the workers
	service, _ := env.Env.GetUserService(ThisServiceName)

	// Pre-compute/lookup the sleep duration for backoff
	PoliteSleep = service.GetParamInt64("polite_sleep")

	logger_level := service.GetParamString("debug_level")
	if logger_level != "debug" {
		retryableClient.Logger = nil
	}

	Gateway = NewHostGateway(time.Duration(PoliteSleep) * time.Second)

	go InfoFetchCount()
	WDB = work_db.NewGuestbookDB()

	go queueing.Enqueue(ChQSHP)
	go queueing.ClearCompletedPeriodically()

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)

}
