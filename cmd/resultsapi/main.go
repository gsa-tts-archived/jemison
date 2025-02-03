package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/postgres"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Databases sync.Map
var ChQSHP = make(chan queueing.QSHP)
var ThisServiceName = "resultsapi"
var JDB *postgres.JemisonDB

func setupQueues() {
	env.InitGlobalEnv(ThisServiceName)

	InitializeQueues()

	go queueing.Enqueue(ChQSHP)
}

func setUpEngine() *gin.Engine {
	engine := gin.Default()

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
	}

	return engine
}

func main() {
	env.InitGlobalEnv(ThisServiceName)
	setupQueues()

	JDB = postgres.NewJemisonDB()

	fmt.Println(ThisServiceName, " environment initialized")

	engine := setUpEngine()
	zap.L().Info("listening from resultsapi",
		zap.String("port", env.Env.Port))

	// Local and Cloud should both get this from the environment.
	//nolint:gosec
	err := http.ListenAndServe(":"+env.Env.Port, engine)
	if err != nil {
		zap.L().Error("could not launch HTTP server", zap.Error(err))
	}
}
