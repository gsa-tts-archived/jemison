package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/postgres"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/searching"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Databases       sync.Map
	ChQSHP          = make(chan queueing.QSHP)
	ThisServiceName = "resultsapi"
	JDB             *postgres.JemisonDB
)

// ////////// Setup.
func setUpEngine(staticFilesPath string, templateFilesPath string) *gin.Engine {
	engine := gin.Default()

	// Delete when no longer using ui for debugging.
	engine.StaticFS("/static", gin.Dir(staticFilesPath, true))
	engine.LoadHTMLGlob(templateFilesPath + "/*")

	engine.GET("/:search", func(c *gin.Context) {
		requiredQueryParams, optionalQueryParams := getQueryParams(c)

		zap.L().Info("Query Data: ",
			zap.String("affiliate", requiredQueryParams.Affiliate),
			zap.String("query", requiredQueryParams.SearchQuery))

		res := searching.DoTheSearch(requiredQueryParams.Affiliate, requiredQueryParams.SearchQuery, JDB)
		marshalledResults := searching.ParseTheResults(res, requiredQueryParams, optionalQueryParams)
		queryResults := searching.JSONToStruct(marshalledResults)

		c.IndentedJSON(http.StatusOK, queryResults)
	})

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
	}

	return engine
}

func getQueryParams(c *gin.Context) (searching.RequiredQueryParameters, searching.OptionalQueryParameters) {
	// required query parameters
	var requiredQueryParas searching.RequiredQueryParameters
	requiredQueryParas.Affiliate = c.Query("affiliate")
	requiredQueryParas.SearchQuery = c.Query("query")

	// optional query parameters
	var optionalQueryParams searching.OptionalQueryParameters

	enableHighlighting, err := strconv.ParseBool(c.Query("enable_highlighting"))
	if err == nil {
		optionalQueryParams.EnableHighlighting = enableHighlighting
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err == nil {
		optionalQueryParams.Offset = offset
	}

	sortBy, err := strconv.Atoi(c.Query("sort_by"))
	if err == nil {
		optionalQueryParams.SortBy = sortBy
	}

	sitelimit, err := strconv.Atoi(c.Query("sitelimit"))
	if err == nil {
		optionalQueryParams.Sitelimit = sitelimit
	}

	return requiredQueryParas, optionalQueryParams
}

////////////////////

func main() {
	env.InitGlobalEnv(ThisServiceName)

	s, _ := env.Env.GetUserService(ThisServiceName)
	templateFilesPath := s.GetParamString("template_files_path")
	staticFilesPath := s.GetParamString("static_files_path")

	JDB = postgres.NewJemisonDB()

	zap.L().Info("environment initialized: " + ThisServiceName)

	engine := setUpEngine(staticFilesPath, templateFilesPath)

	zap.L().Info("listening from resultsapi",
		zap.String("port", env.Env.Port))

	// Local and Cloud should both get this from the environment.
	//nolint:gosec
	err := http.ListenAndServe(":"+env.Env.Port, engine)
	if err != nil {
		zap.L().Error("could not launch HTTP server", zap.Error(err))
	}
}
