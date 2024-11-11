package main

import (
	"log"
	"net/http"
	"os"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ThisServiceName = "admin"

type FetchRequestInput struct {
	Scheme string `json:"scheme" maxLength:"10" doc:"Resource scheme"`
	Host   string `json:"host" maxLength:"500" doc:"Host of resource"`
	Path   string `json:"path" maxLength:"1500" doc:"Path to resource"`
	ApiKey string `json:"api-key"`
}

// https://dev.to/kashifsoofi/rest-api-with-go-chi-and-inmemory-store-43ag
func CrawlRequestHandler(c *gin.Context) {
	var fri FetchRequestInput
	if err := c.BindJSON(&fri); err != nil {
		return
	}

	if fri.ApiKey == os.Getenv("API_KEY") || true {
		zap.L().Debug("api enqueue", zap.String("host", fri.Host), zap.String("path", fri.Path))

		queueing.InsertFetch(fri.Scheme, fri.Host, fri.Path)

		c.IndentedJSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func main() {
	env.InitGlobalEnv(ThisServiceName)
	InitializeQueues()

	engine := common.InitializeAPI()

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
		v1.PUT("/crawl", CrawlRequestHandler)
	}

	log.Println("environment initialized")

	// // Init a cache for the workers
	// service, _ := env.Env.GetUserService("admin")

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)

}
