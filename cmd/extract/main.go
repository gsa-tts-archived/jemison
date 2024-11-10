package main

import (
	"log"
	"net/http"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"go.uber.org/zap"
)

var ThisServiceName = "extract"

func main() {
	env.InitGlobalEnv()
	InitializeQueues()
	log.Println("environment initialized")

	routers := common.InitializeAPI()

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, routers)

}
