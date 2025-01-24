package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type requiredQueryParameters struct {
	affiliate   string
	searchQuery string
}

type optionalQueryParameters struct {
	enableHighlighting bool
	offset             int
	sortBy             int
	sitelimit          int
}

type QueryWebResultsData struct {
	Title           string `json:"title"`
	URL             string `json:"url"`
	Snippet         string `json:"snippet"`
	PublicationDate string `json:"publication_date"`
	ThumbnailURL    string `json:"thumbnail_url"`
}

type QueryWebData struct {
	Total              int                   `json:"total"`
	NextOffset         int                   `json:"next_offset"`
	SpellingCorrection string                `json:"spelling_correction"`
	Results            []QueryWebResultsData `json:"results"`
}

type QueryData struct {
	SearchedQuery            string       `json:"query"`
	Web                      QueryWebData `json:"web"`
	TextBestBets             []string     `json:"text_best_bets"`
	GraphicBestBets          []string     `json:"graphic_best_bets"`
	HealthTopics             []string     `json:"health_topics"`
	JobOpenings              []string     `json:"job_openings"`
	FederalRegisterDocuments []string     `json:"federal_register_documents"`
	RelatedSearchTerms       []string     `json:"related_search_terms"`
}

func setupQueues() {
	env.InitGlobalEnv(ThisServiceName)

	InitializeQueues()

	go queueing.Enqueue(ChQSHP)
}

func getQueryParams(c *gin.Context) {
	// required query parameters
	var requiredQueryParas requiredQueryParameters
	requiredQueryParas.affiliate = c.Query("affiliate")
	requiredQueryParas.searchQuery = c.Query("query")
	fmt.Println("affiliate: ", requiredQueryParas.affiliate, " query: ", requiredQueryParas.searchQuery)

	// optional query parameters
	var optionalQueryParams optionalQueryParameters

	enableHighlighting, err := strconv.ParseBool(c.Query("enable_highlighting"))
	if err == nil {
		optionalQueryParams.enableHighlighting = enableHighlighting
		fmt.Println("enableHighlighting: ", optionalQueryParams.enableHighlighting)
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err == nil {
		optionalQueryParams.offset = offset
		fmt.Println("offset: ", optionalQueryParams.offset)
	}

	sortBy, err := strconv.Atoi(c.Query("sort_by"))
	if err == nil {
		optionalQueryParams.sortBy = sortBy
		fmt.Println("sortBy: ", optionalQueryParams.sortBy)
	}

	sitelimit, err := strconv.Atoi(c.Query("sitelimit"))
	if err == nil {
		optionalQueryParams.sitelimit = sitelimit
		fmt.Println("sitelimit: ", optionalQueryParams.sitelimit)
	}
}

// ///// Contained Functions are used to test and create data for JSON to present to Client ////
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func readHttpContent() string {
	//note this is already formatted correctly but we will still grab the data format it into the struct we need and rebound as a JSON to serve to our client
	var url = os.Getenv("TMPURL")

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)
	req.Header.Set("User-Agent", "")
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	fmt.Printf("Response type:%T\n", resp)
	bytes, err := io.ReadAll(resp.Body)
	checkError(err)
	content := string(bytes)
	return content
}

func makeStructFromJSON(content string) []QueryData {
	data := make([]QueryData, 0)
	decoder := json.NewDecoder(strings.NewReader(content))
	for {
		var result QueryData
		err := decoder.Decode(&result)
		if err != nil {
			if err == io.EOF {
				break // Reached end of the stream
			}
			panic(err)
		}
		data = append(data, result)
	}
	return data
}

/////////////////////////////////////////////////////////////////////////////////////

func getQueryResponse() []QueryData {
	// TODO: get data from database
	// using "dummy" site to get data from
	content := readHttpContent()
	data := makeStructFromJSON(content)
	return data
}

func makeJSONFromStruct(data []QueryData) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return jsonData
}

func returnQueryResponse(data []QueryData) {
	jsonData := makeJSONFromStruct(data)
	fmt.Printf("%s", string(jsonData))
	// TODO: return the JSON to client
}

func setUpEngine(staticFilesPath string, templateFilesPath string) *gin.Engine {
	engine := gin.Default()

	// will we need the two instructions below? I think not because there will be no ui
	engine.StaticFS("/static", gin.Dir(staticFilesPath, true))
	engine.LoadHTMLGlob(templateFilesPath + "/*")

	engine.GET("/:search", func(c *gin.Context) {
		getQueryParams(c)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	// simulate parsing data from database
	data := getQueryResponse()

	// return the data as a JSON object
	returnQueryResponse(data)

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

	fmt.Println(ThisServiceName, " environment initialized")

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
