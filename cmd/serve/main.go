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

var Databases sync.Map //map[string]*sql.DB
var ChQSHP = make(chan queueing.QSHP)
var ThisServiceName = "serve"
var JDB *postgres.JemisonDB

func main() {
	env.InitGlobalEnv(ThisServiceName)
	//s3 := kv.NewS3(ThisServiceName)
	InitializeQueues()

	go queueing.Enqueue(ChQSHP)

	s, _ := env.Env.GetUserService(ThisServiceName)
	template_files_path := s.GetParamString("template_files_path")
	static_files_path := s.GetParamString("static_files_path")

	external_host := s.GetParamString("external_host")
	external_port := s.GetParamInt64("external_port")

	JDB = postgres.NewJemisonDB()

	log.Println("environment initialized")

	zap.L().Info("serve environment",
		zap.String("template_files_path", template_files_path),
		zap.String("external_host", external_host),
		zap.Int64("external_port", external_port),
	)

	/////////////////////
	// Server/API
	engine := gin.Default()
	// engine.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusMovedPermanently, "/search/"+start)
	// })
	// engine.GET("/search", func(c *gin.Context) {
	// 	c.Redirect(http.StatusMovedPermanently, "/search/"+start)
	// })
	engine.StaticFS("/static", gin.Dir(static_files_path, true))
	// engine.GET("/search/:host", ServeHost)

	engine.LoadHTMLGlob(template_files_path + "/*")

	base_params := gin.H{
		"scheme":      "http",
		"search_host": "localhost",
		"search_port": "10000",
	}

	engine.GET("/:tld/:domain", func(c *gin.Context) {
		base_params["tld"] = c.Param("tld")
		base_params["domain"] = c.Param("domain")
		base_params["fqdn"] = c.Param("domain") + "." + c.Param("tld")
		c.HTML(http.StatusOK, "index.tmpl", base_params)
	})

	engine.GET("/:tld/:domain/:subdomain", func(c *gin.Context) {
		base_params["tld"] = c.Param("tld")
		base_params["domain"] = c.Param("domain")
		base_params["subdomain"] = c.Param("subdomain")
		base_params["fqdn"] = c.Param("subdomain") + "." + c.Param("domain") + "." + c.Param("tld")
		c.HTML(http.StatusOK, "index.tmpl", base_params)
	})

	v1 := engine.Group("/api")
	{
		v1.GET("/heartbeat", common.Heartbeat)
		v1.POST("/search", SearchHandler)
	}

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)

}
