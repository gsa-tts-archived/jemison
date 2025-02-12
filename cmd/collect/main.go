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

var ThisServiceName = "collect"

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
	// Initialize environment for "collect" service
	env.InitGlobalEnv(ThisServiceName)

	// Initialize JSON schemas
	if err := InitializeSchemas(); err != nil {
		zap.L().Fatal("failed to initialize schemas", zap.Error(err))
	}

	// Initialize worker queues
	setupQueues()

	// Create database connection
	JDB = postgres.NewJemisonDB()

	fmt.Println(ThisServiceName, " environment initialized")

	// Setting up HTTP engine
	engine := setUpEngine()

	zap.L().Info("listening from collect", zap.String("port", env.Env.Port))

	// Start the service
	//
	//nolint:gosec // G114: Ignoring timeout settings for demonstration purposes or due to intentional design
	if err := http.ListenAndServe(":"+env.Env.Port, engine); err != nil {
		zap.Error(err)
	}
}
