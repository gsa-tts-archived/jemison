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
func FetchRequestHandler(c *gin.Context) {
	var fri FetchRequestInput
	if err := c.BindJSON(&fri); err != nil {
		return
	}
	if fri.ApiKey == os.Getenv("API_KEY") || true {
		zap.L().Debug("fetch enqueue", zap.String("host", fri.Host), zap.String("path", fri.Path))
		queueing.InsertFetch(fri.Scheme, fri.Host, fri.Path)
		c.IndentedJSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func EntreeRequestHandler(c *gin.Context) {
	var fri FetchRequestInput
	full := c.Param("fullorone")
	hallPass := c.Param("hallpass")

	if err := c.BindJSON(&fri); err != nil {
		return
	}
	if fri.ApiKey == os.Getenv("API_KEY") || true {
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
			zap.Bool("full", fullB),
			zap.Bool("hallpass", hallPassB))

		queueing.InsertEntree(fri.Scheme, fri.Host, fri.Path, fullB, hallPassB)

		c.IndentedJSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

// func JobCountHandler(c *gin.Context) {
// 	ctx := context.Background()
// 	db_url, _ := env.Env.GetDatabaseUrl(env.QueueDatabase)

// 	conn, err := pgx.Connect(ctx, db_url)
// 	if err != nil {
// 	}
// 	defer conn.Close(ctx)

// 	queries := pg_schemas.New(conn)
// 	counts := make(map[string]map[string]int64)
// 	for _, service := range []string{"fetch", "extract", "pack", "validate", "walk", "serve"} {
// 		service_counts := make(map[string]int64)
// 		for _, state := range []string{"completed", "running", "retryable"} {
// 			count, _ := queries.CountJobs(ctx, pg_schemas.CountJobsParams{
// 				Kind:  service,
// 				State: state,
// 			})
// 			service_counts[state] = count
// 		}
// 		counts[service] = service_counts
// 	}

// 	c.IndentedJSON(http.StatusOK, gin.H{
// 		"status": "ok",
// 		"counts": counts,
// 	})
// }

func main() {
	env.InitGlobalEnv(ThisServiceName)
	InitializeQueues()

	engine := common.InitializeAPI()

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
		v1.PUT("/fetch", FetchRequestHandler)
		v1.PUT("/entree/:fullorone/:hallpass", EntreeRequestHandler)
		// v1.GET("/jobs", JobCountHandler)
	}

	log.Println("environment initialized")

	// // Init a cache for the workers
	// service, _ := env.Env.GetUserService("admin")

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)

}
