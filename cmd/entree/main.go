package main

import (
	"log"
	"net/http"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"go.uber.org/zap"
)

// var expirable_cache expirable.Cache[string, int]

var ThisServiceName = "entree"
var HostIdMap = make(map[string]int64)

func main() {
	env.InitGlobalEnv(ThisServiceName)
	InitializeQueues()
	// service, _ := env.Env.GetUserService(ThisServiceName)

	engine := common.InitializeAPI()

	log.Println("environment initialized")

	HostIdMap = upsertUniqueHosts()
	crontab()

	WDB = NewGuestbookDB()

	// // For a short period of time, we need to not re-insert things while crawling.
	// // If we don't actually fetch a page, we can end up re-queueing it. Because
	// ttl := service.GetParamInt64("cache-ttl")
	// expirable_cache = expirable.NewCache[string, int]().WithTTL(time.Duration(ttl) * time.Second)

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)

}
