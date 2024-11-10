package main

import (
	"log"
	"net/http"
	"time"

	expirable "github.com/go-pkgz/expirable-cache/v3"
	"github.com/patrickmn/go-cache"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"go.uber.org/zap"
)

var expirable_cache expirable.Cache[string, int]
var RecentlyVisitedCache *cache.Cache

func main() {

	env.InitGlobalEnv()
	InitializeQueues()

	log.Println("environment initialized")
	service, _ := env.Env.GetUserService("walk")

	engine := common.InitializeAPI()

	ttl := service.GetParamInt64("cache-ttl")
	expirable_cache = expirable.NewCache[string, int]().WithTTL(time.Duration(ttl) * time.Second)

	RecentlyVisitedCache = cache.New(
		time.Duration(service.GetParamInt64("polite_cache_default_expiration"))*time.Second,
		time.Duration(service.GetParamInt64("polite_cache_cleanup_interval"))*time.Second)

	log.Println("WALK listening on", env.Env.Port)

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)
}
