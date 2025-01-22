package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/postgres"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Databases sync.Map

var ChQSHP = make(chan queueing.QSHP)

var ThisServiceName = "resultsapi"

var JDB *postgres.JemisonDB

type fakeResult struct {
	SearchedQuery      string `json:"searchedQuery"`
	Title              string `json:"title"`
	Snippet            string `json:"snippet"`
	PublicationDate    string `json:"publication_date"`
	ThumbnailUrl       string `json:"thumbnail_url"`
	TextBestBets       string `json:"text_best_bets"`
	GraphicBestBets    string `json:"graphic_best_bets"`
	HealthTopics       string `json:"health_yopics"`
	JobOpenings        string `json:"job_openings"`
	RelatedSearchTerms string `json:"related_search_terms"`
}

fakeResults := fakeResult{
	{SearchedQuery: "https://www.nasa.gov/news-release/nasa-releases-detailed-global-climate-change-projections", Title: "NASA Releases Detailed Global Climate Change Projections", Snippet: "NASA has released data showing how temperature and rainfall patterns worldwide may change through the year 2100 because of growing concentrations of greenhouse gases in Earth’s atmosphere.", PublicationDate: "Jun 09, 2015", ThumbnailUrl: "https://www.nasa.gov/wp-content/uploads/2015/06/15-115.jpg?resize=2000,935", TextBestBets: " ", GraphicBestBets: " ", HealthTopics: " ", JobOpenings: " ", RelatedSearchTerms: " "},
}

func setupQueues() {
	env.InitGlobalEnv(ThisServiceName)

	InitializeQueues()

	go queueing.Enqueue(ChQSHP)
}

func setUpEngine(staticFilesPath string, templateFilesPath string) *gin.Engine {
	engine := gin.Default()

	// will we need the two instructions below? I think not because there will be no ui
	engine.StaticFS("/static", gin.Dir(staticFilesPath, true))
	engine.LoadHTMLGlob(templateFilesPath + "/*")

	engine.GET("/:search", func(c *gin.Context) {
		//required query parameters
		affiliate := c.Query("affiliate")
		searchQuery := c.Query("query")
		log.Println("affiliate: ", affiliate, " query: ", searchQuery)

		//optional query parameters
		// enable_highlighting := c.Query("enable_highlighting")
		// offset := c.Query("offset")
		// sort_by := c.Query("sort_by")
		// sitelimit := c.Query("sitelimit")

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"affiliate": affiliate,
			"query":     searchQuery,
			"faker":     fakeResults,
		})
	})

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
	}

	return engine
}

//nolint:funlen
func main() {
	setupQueues()

	s, _ := env.Env.GetUserService(ThisServiceName)
	templateFilesPath := s.GetParamString("template_files_path")
	staticFilesPath := s.GetParamString("static_files_path")

	externalHost := s.GetParamString("external_host")
	externalPort := s.GetParamInt64("external_port")

	JDB = postgres.NewJemisonDB()

	log.Println(ThisServiceName, " environment initialized")

	zap.L().Info("resultsapi environment",
		zap.String("template_files_path", templateFilesPath),
		zap.String("external_host", externalHost),
		zap.Int64("external_port", externalPort),
	)

	engine := setUpEngine(staticFilesPath, templateFilesPath)

	zap.L().Info("listening from resultsapi",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	//nolint:gosec
	err := http.ListenAndServe(":"+env.Env.Port, engine)
	if err != nil {
		zap.Error(err)
	}
}
