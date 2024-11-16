package main

import (
	"log"
	"net/http"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"go.uber.org/zap"
)

var ThisServiceName = "entree"
var HostIdMap = make(map[string]int64)

func main() {
	env.InitGlobalEnv(ThisServiceName)
	InitializeQueues()

	engine := common.InitializeAPI()

	log.Println("environment initialized")

	HostIdMap = upsertUniqueHosts()
	crontab()

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)

}
