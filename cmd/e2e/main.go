package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/postgres"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ThisServiceName = "e2e"
var ChQSHP = make(chan queueing.QSHP)

var JDB *postgres.JemisonDB

func main() {
	env.InitGlobalEnv(ThisServiceName)
	s, _ := env.Env.GetUserService(ThisServiceName)

	// For the heartbeat
	engine := common.InitializeAPI()

	template_files_path := s.GetParamString("template_files_path")
	static_files_path := s.GetParamString("static_files_path")

	// For the index template
	// This provides links to the HTML files
	// in the static directory.
	engine.LoadHTMLGlob(template_files_path + "/*")

	base_params := gin.H{
		"scheme":      "http",
		"search_host": "e2e.gov",
		"search_port": env.Env.Port,
	}

	files, err := os.ReadDir(static_files_path)
	if err != nil {
		log.Fatal(err)
	}

	filenames := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}

	engine.GET("/www", func(c *gin.Context) {
		// d64_start, _ := strconv.ParseInt("010000DD00000000", 16, 64)
		// d64_end, _ := strconv.ParseInt("010000DDFFFFFF00", 16, 64)
		// base_params["d64_start"] = d64_start
		// base_params["d64_end"] = d64_end
		base_params["filenames"] = filenames
		c.HTML(http.StatusOK, "index.tmpl", base_params)
	})

	engine.HEAD("/www", func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		base_params["filenames"] = filenames

		c.HTML(http.StatusOK, "index.tmpl", base_params)
	})

	go queueing.Enqueue(ChQSHP)

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))

	engine.StaticFS("/static", gin.Dir(static_files_path, true))

	go http.ListenAndServe(":80", engine)
	http.ListenAndServe(":"+env.Env.Port, engine)
}
