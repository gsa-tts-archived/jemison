package main

import (
	"log"
	"net/http"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

var polite_sleep int64
var ThisServiceName = "fetch"

var RetryClient *http.Client

var ch_qshp = make(chan queueing.QSHP)

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
	polite_sleep = service.GetParamInt64("polite_sleep")

	go InfoFetchCount()
	go queueing.Enqueue(ch_qshp)

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)

}
