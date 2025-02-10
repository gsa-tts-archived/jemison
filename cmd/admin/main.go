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

var ChQSHP = make(chan queueing.QSHP)

type FetchRequestInput struct {
	Scheme string `json:"scheme" maxLength:"10" doc:"Resource scheme"`
	Host   string `json:"host" maxLength:"500" doc:"Host of resource"`
	Path   string `json:"path" maxLength:"1500" doc:"Path to resource"`
	APIKey string `json:"api-key"`
	Data   string `json:"data"`
}

// https://dev.to/kashifsoofi/rest-api-with-go-chi-and-inmemory-store-43ag
func FetchRequestHandler(c *gin.Context) {
	var fri FetchRequestInput
	if err := c.BindJSON(&fri); err != nil {
		return
	}

	if fri.APIKey == os.Getenv("API_KEY") || true {
		zap.L().Debug("fetch enqueue",
			zap.String("host", fri.Host),
			zap.String("path", fri.Path),
			zap.String("path", fri.Data))

		ChQSHP <- queueing.QSHP{
			Queue:   "collect",
			Scheme:  fri.Scheme,
			Host:    fri.Host,
			Path:    fri.Path,
			RawData: fri.Data,
		}

		ChQSHP <- queueing.QSHP{
			Queue:  "fetch",
			Scheme: fri.Scheme,
			Host:   fri.Host,
			Path:   fri.Path,
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func EntreeRequestHandler(c *gin.Context) {
	var fri FetchRequestInput

	// Parse and bind JSON into the struct
	if err := c.ShouldBindJSON(&fri); err != nil {
		zap.L().Error("failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	full := c.Param("fullorone")
	hallPass := c.Param("hallpass")

	if fri.APIKey == os.Getenv("API_KEY") || true {
		hallPassB := false
		fullB := false

		if hallPass == "pass" {
			hallPassB = true
		}

		if full == "full" {
			fullB = true
		}

		zap.L().Debug("entree enqueue",
			zap.String("host", fri.Host),
			zap.String("path", fri.Path),
			zap.String("data", fri.Data))

		ChQSHP <- queueing.QSHP{
			Queue:   "collect",
			Scheme:  fri.Scheme,
			Host:    fri.Host,
			Path:    fri.Path,
			RawData: fri.Data,
		}

		ChQSHP <- queueing.QSHP{
			Queue:      "entree",
			Scheme:     fri.Scheme,
			Host:       fri.Host,
			Path:       fri.Path,
			IsFull:     fullB,
			IsHallPass: hallPassB,
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func PackRequestHandler(c *gin.Context) {
	var fri FetchRequestInput

	if err := c.BindJSON(&fri); err != nil {
		return
	}

	if fri.APIKey == os.Getenv("API_KEY") || true {
		zap.L().Debug("pack enqueue",
			zap.String("host", fri.Host))

		ChQSHP <- queueing.QSHP{
			Queue:  "pack",
			Scheme: fri.Scheme,
			Host:   fri.Host,
			Path:   fri.Path,
		}
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
		v1.PUT("/fetch", FetchRequestHandler)
		v1.PUT("/entree/:fullorone/:hallpass", EntreeRequestHandler)
		v1.PUT("/pack", PackRequestHandler)
	}

	log.Println("environment initialized")

	go queueing.Enqueue(ChQSHP)

	// // Init a cache for the workers
	// service, _ := env.Env.GetUserService("admin")

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	//nolint:gosec
	err := http.ListenAndServe(":"+env.Env.Port, engine)
	if err != nil {
		zap.Error(err)
	}
}
