package main

import (
	"net/http"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/postgres"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ThisServiceName = "e2e"
var ChQSHP = make(chan queueing.QSHP)

var JDB *postgres.JemisonDB

func main() {
	env.InitGlobalEnv(ThisServiceName)
	s, _ := env.Env.GetUserService(ThisServiceName)

	engine := common.InitializeAPI()

	go queueing.Enqueue(ChQSHP)

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))

	static_files_path := s.GetParamString("static_files_path")
	engine.StaticFS("/static", gin.Dir(static_files_path, true))

	http.ListenAndServe(":"+env.Env.Port, engine)
}
