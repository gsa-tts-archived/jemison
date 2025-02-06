package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/GSA-TTS/jemison/config"
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

// ////////// Setup //////////
func setupQueues() {
	env.InitGlobalEnv(ThisServiceName)

	InitializeQueues()

	go queueing.Enqueue(ChQSHP)
}

func setUpEngine(staticFilesPath string, templateFilesPath string) *gin.Engine {
	engine := gin.Default()

	// TODO: Delete when no longer using ui for debugging
	engine.StaticFS("/static", gin.Dir(staticFilesPath, true))
	engine.LoadHTMLGlob(templateFilesPath + "/*")

	engine.GET("/:search", func(c *gin.Context) {
		//required query parameters
		// affiliate := c.Query("affiliate")
		// searchQuery := c.Query("query")
		requiredQueryParams, optionalQueryParams := getQueryParams(c)

		zap.L().Info("Query Data: ",
			zap.String("affiliate", requiredQueryParams.affiliate),
			zap.String("query", requiredQueryParams.searchQuery))

		res := doTheSearch(requiredQueryParams.affiliate, requiredQueryParams.searchQuery)
		pretty_res := parseTheResults(res, requiredQueryParams, optionalQueryParams)
		//optional query parameters
		// enable_highlighting := c.Query("enable_highlighting")
		// offset := c.Query("offset")
		// sort_by := c.Query("sort_by")
		// sitelimit := c.Query("sitelimit")

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"res":        res,
			"pretty_res": pretty_res,
		})
	})

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
	}

	return engine
}

func getQueryParams(c *gin.Context) (requiredQueryParameters, optionalQueryParameters) {
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
	return requiredQueryParas, optionalQueryParams
}

////////////////////

// ////////// Searching //////////
func doTheSearch(affiliate string, searchQuery string) []SearchResult {
	domain64Start, domain64End := getD64(affiliate + ".gov")

	sri := SearchRequestInput{
		Host:          affiliate + ".gov",
		Path:          "",
		Terms:         searchQuery,
		Domain64Start: domain64Start,
		Domain64End:   domain64End,
	}

	rows, duration, err := runQuery(sri)

	zap.L().Info("Queried Answer:",
		zap.Any("rows: ", rows),
		zap.Any("duration", duration),
		zap.Any("err", err))

	return rows
}

func getD64(affiliate string) (string, string) {
	var subdomain, domain, tld string

	subdomain, domain, tld = parseAffiliate(affiliate)

	var d64Start, d64End int64

	// top level domain
	d64Start, _ = strconv.ParseInt(fmt.Sprintf("%02x00000000000000", tld), 16, 64)
	d64End, _ = strconv.ParseInt(fmt.Sprintf("%02xFFFFFFFFFFFF00", tld), 16, 64)

	// domain
	if domain != "" {
		start := config.RDomainToDomain64(fmt.Sprintf("%s.%s", tld, domain))
		d64Start, _ = strconv.ParseInt(fmt.Sprintf("%s00000000", start), 16, 64)
		d64End, _ = strconv.ParseInt(fmt.Sprintf("%sFFFFFF00", start), 16, 64)
	} else {
		sD64Start := fmt.Sprintf("%d", d64Start)
		sD64End := fmt.Sprintf("%d", d64End)

		return sD64Start, sD64End
	}

	//subdomain
	if subdomain != "" {
		fqdn := fmt.Sprintf("%s.%s.%s", subdomain, domain, tld)
		start, _ := config.FQDNToDomain64(fqdn)
		d64Start = start
		d64End = start + 1
	}

	sD64Start := fmt.Sprintf("%d", d64Start)
	sD64End := fmt.Sprintf("%d", d64End)

	return sD64Start, sD64End
}

func parseAffiliate(affiliate string) (string, string, string) {
	tld := ""
	domain := ""
	subdomain := ""
	delimiter := "."

	results := strings.Split(affiliate, delimiter)

	if len(results) == 3 {
		subdomain = results[0]
		domain = results[1]
		tld = results[2]
	} else if len(results) == 2 {
		domain = results[0]
		tld = results[1]
	} else {
		tld = results[0]
	}

	return subdomain, domain, tld
}

////////////////////

// ////////// Returning Results //////////
func parseTheResults(results []SearchResult, requiredQueryParams requiredQueryParameters, optionalQueryParams optionalQueryParameters) string {

	//create array of results {JSONResults}
	jSONResults := createJSONResults(results)

	//create webJSON
	webResults := createWebResults(jSONResults, optionalQueryParams)

	//create wholeJSON
	wholeJSON := createWholeJSON(webResults, requiredQueryParams)

	return wholeJSON
}

func createJSONResults(results []SearchResult) []QueryWebResultsData {
	// var JSONResults []string
	var JSONResults []QueryWebResultsData
	for _, r := range results {
		//convert searchresult into a json object that matches current resultAPI
		jsonStr, err := getJSONString(r)
		qwrd := createQueryWebResultsData(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("NEW ENTRY: " + r.PageTitle + "NEW ENTRY JSON: " + jsonStr)
		//append to JSONResults
		JSONResults = append(JSONResults, qwrd)
	}
	return JSONResults
}

func createWebResults(jSONResults []QueryWebResultsData, optionalQueryParams optionalQueryParameters) QueryWebData {
	total := 5
	nextOffset := optionalQueryParams.offset
	spellingCorrections := "null"
	strc := QueryWebData{total, nextOffset, spellingCorrections, jSONResults}

	return strc
}

func createWholeJSON(webResults QueryWebData, requiredQueryParams requiredQueryParameters) string {
	query := requiredQueryParams.searchQuery
	var tretBestBets []string
	var graphicBestBets []string
	var healthTopics []string
	var jobOpenings []string
	var federalRegisterDocuments []string
	var relatedSearchTerms []string

	strc := QueryData{query, webResults, tretBestBets, graphicBestBets, healthTopics, jobOpenings, federalRegisterDocuments, relatedSearchTerms}
	data := structToJSON(strc)
	return string(data)
}

func getJSONString(strc interface{}) (string, error) {
	//? can I convert a struct to another struct?

	//convert struct to JSON
	data := structToJSON(strc)

	//from JSON convert to new struct
	var searchResultJSON SearchResultJSON
	json.Unmarshal([]byte(data), &searchResultJSON)

	//TODO: get these values
	searchResultJSON.PublicationDate = "null"
	searchResultJSON.ThumbnailUrl = "null"

	//convert new struct back to JSON
	j_data := structToJSON(searchResultJSON)

	return string(j_data), nil
}

func createQueryWebResultsData(strc interface{}) QueryWebResultsData {
	//? can I convert a struct to another struct?

	//convert struct to JSON
	data := structToJSON(strc)

	//from JSON convert to new struct
	var searchResultJSON QueryWebResultsData
	json.Unmarshal([]byte(data), &searchResultJSON)
	return searchResultJSON

}

func structToJSON(strc interface{}) []byte {
	data, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

////////////////////

func main() {
	env.InitGlobalEnv(ThisServiceName)
	setupQueues()

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
