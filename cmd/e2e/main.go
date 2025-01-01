package main

import (
	"net/http"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/postgres"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"go.uber.org/zap"
)

var ThisServiceName = "e2e"
var ChQSHP = make(chan queueing.QSHP)

var JDB *postgres.JemisonDB

func main() {
	env.InitGlobalEnv(ThisServiceName)

	engine := common.InitializeAPI()

	go queueing.Enqueue(ChQSHP)

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))

	// Serve the static directory at /static
	engine.Static("/static", "./static")
	http.ListenAndServe(":"+env.Env.Port, engine)
}
