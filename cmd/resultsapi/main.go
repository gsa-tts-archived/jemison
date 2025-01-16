package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/postgres"
	"github.com/GSA-TTS/jemison/internal/postgres/search_db"
	"github.com/GSA-TTS/jemison/internal/postgres/work_db"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Databases sync.Map

var ChQSHP = make(chan queueing.QSHP)

var ThisServiceName = "resultsapi"

var JDB *postgres.JemisonDB

type fakeResult struct {
	URL             string `json:"url"`
	Title           string `json:"title"`
	Snippet         string `json:"snippet"`
	PublicationDate string `json:"publication_date"`
	ThumbnailUrl    string `json:"thumbnail_url"`
}

var fakeResults = []fakeResult{
	{URL: "https://www.nasa.gov/news-release/nasa-releases-detailed-global-climate-change-projections", Title: "NASA Releases Detailed Global Climate Change Projections", Snippet: "NASA has released data showing how temperature and rainfall patterns worldwide may change through the year 2100 because of growing concentrations of greenhouse gases in Earth’s atmosphere.", PublicationDate: "Jun 09, 2015", ThumbnailUrl: "https://www.nasa.gov/wp-content/uploads/2015/06/15-115.jpg?resize=2000,935"},
}

func addMetadata(m map[string]any) map[string]any {
	pathCount, err := JDB.WorkDBQueries.PathsInDomain64Range(context.Background(),
		work_db.PathsInDomain64RangeParams{
			D64Start: m["d64_start"].(int64),
			D64End:   m["d64_end"].(int64),
		})
	if err != nil {
		zap.L().Error(err.Error())

		pathCount = 0
	}

	m["pageCount"] = pathCount

	bodyCount, err := JDB.SearchDBQueries.BodiesInDomain64Range(context.Background(),
		search_db.BodiesInDomain64RangeParams{
			D64Start: m["d64_start"].(int64),
			D64End:   m["d64_end"].(int64),
		})
	if err != nil {
		zap.L().Error(err.Error())

		bodyCount = 0
	}

	m["bodyCount"] = bodyCount

	return m
}

//nolint:funlen
func main() {
	env.InitGlobalEnv(ThisServiceName)

	InitializeQueues()

	go queueing.Enqueue(ChQSHP)

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

	/////////////////////
	/// Results API ////
	engine := gin.Default()

	// will we need the two instructions below? I think not because there will be no ui
	engine.StaticFS("/static", gin.Dir(staticFilesPath, true))
	engine.LoadHTMLGlob(templateFilesPath + "/*")

	// baseParams := gin.H{
	// 	"scheme":      "http",
	// 	"search_host": "localhost",
	// 	"search_port": "10008",
	// }

	engine.GET("/:search", func(c *gin.Context) {
		affiliate := c.Query("affiliate")
		query := c.Query("query")
		log.Println("affiliate: ", affiliate, " query: ", query)

		// why are we doing all of this and do we need it?
		// tld := config.GetTLD(c.Param("search"))
		// d64Start, _ := strconv.ParseInt(fmt.Sprintf("%02x00000000000000", tld), 16, 64)
		// d64End, _ := strconv.ParseInt(fmt.Sprintf("%02xFFFFFFFFFFFF00", tld), 16, 64)
		// baseParams["tld"] = c.Param("tld")
		// delete(baseParams, "domain")
		// delete(baseParams, "subdomain")
		// baseParams["fqdn"] = c.Param("tld")
		// baseParams["d64_start"] = d64Start
		// baseParams["d64_end"] = d64End
		// baseParams = addMetadata(baseParams)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"affiliate": affiliate,
			"query":     query,
			"faker":     fakeResults,
		})
	})

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
	}

	zap.L().Info("listening from resultsapi",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	//nolint:gosec
	err := http.ListenAndServe(":"+env.Env.Port, engine)
	if err != nil {
		zap.Error(err)
	}
}
