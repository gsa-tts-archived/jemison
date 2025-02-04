package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

type PostBody struct {
	terms     string `json:"terms"`
	host      string `json:"host"`
	d64_start int64  `json:"d64_start"`
	d64_end   int64  `json:"d64_end"`
}

func setupQueues() {
	env.InitGlobalEnv(ThisServiceName)

	InitializeQueues()

	go queueing.Enqueue(ChQSHP)
}

func sendPOST(url string, body PostBody) io.ReadCloser {
	fmt.Println("in SENDPOST")
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	req, _ := http.NewRequest("POST", url, payloadBuf)

	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil {
		// return e
	}

	defer res.Body.Close()

	fmt.Println("response Status:", res.Status)
	// Print the body to the stdout
	io.Copy(os.Stdout, res.Body)

	return res.Body
}

func getPostURL(scheme string, host string) string {
	var posturl = scheme + "://www." + host + "/api/search"
	fmt.Println("in getPostURL: ", posturl)
	return posturl
}

func getD64() (int64, int64) {
	return 67, 89
}

func doTheSearch(affiliate string, searchQuery string) {
	fmt.Println("in DOTHESEARCH")
	// Makes a fetch (in javascript) query to post_url for a POSt method call with header & JSON body

	// - post_url is "{{.scheme}}://" + host + "/api/search" ==> http://localhost:10000/api/search
	// - scheme is http from main.go in serve
	// TODO: where is scheme coming from here?
	// - host is window.location.host = domain + post number if it exists
	// TODO: where is host coming from here?
	var post_url = getPostURL("http", "localhost:10000")

	// - JSON body is stringified data
	// 	- data is object of terms, host, d64_start & d_64end
	// TODO get host, d64s
	var postBody PostBody
	postBody.terms = searchQuery
	postBody.host = affiliate
	postBody.d64_start, postBody.d64_end = getD64()
	sendPOST(post_url, postBody)

	// - returns a JSON
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

		doTheSearch(affiliate, searchQuery)
		//optional query parameters
		// enable_highlighting := c.Query("enable_highlighting")
		// offset := c.Query("offset")
		// sort_by := c.Query("sort_by")
		// sitelimit := c.Query("sitelimit")

		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
	}

	return engine
}

func main() {
	env.InitGlobalEnv(ThisServiceName)
	setupQueues()

	s, _ := env.Env.GetUserService(ThisServiceName)
	templateFilesPath := s.GetParamString("template_files_path")
	staticFilesPath := s.GetParamString("static_files_path")

	JDB = postgres.NewJemisonDB()

	fmt.Println(ThisServiceName, " environment initialized")

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
