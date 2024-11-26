package main

import (
	"net/http"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"go.uber.org/zap"
)

var ThisServiceName string = "balance"
var TheSpaceMap *SpaceMap

func main() {
	env.InitGlobalEnv(ThisServiceName)
	InitializeQueues()

	engine := common.InitializeAPI()
	TheSpaceMap = NewSpaceMap()

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
		// v1.PUT("/fetch", FetchRequestHandler)
		// v1.PUT("/entree/:fullorone/:hallpass", EntreeRequestHandler)
	}

	// go queueing.Enqueue(ChQSHP)

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	err := http.ListenAndServe(":"+env.Env.Port, engine)
	if err != nil {
		zap.L().Error("could not take the first step",
			zap.String("port", env.Env.Port))
	}
}
