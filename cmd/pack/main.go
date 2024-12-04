package main

import (
	"log"
	"net/http"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/search_db/search_db"
)

var ThisServiceName = "pack"
var ChFinalize = make(chan string)
var ChQSHP = make(chan queueing.QSHP)
var PHL *PerHostLock = nil
var SDB0 *search_db.SearchDB

func main() {
	env.InitGlobalEnv(ThisServiceName)

	InitializeQueues()
	engine := common.InitializeAPI()

	log.Println("environment initialized")

	PHL = NewPerHostLock()
	SDB0 = search_db.NewSearchDB("s0")

	go FinalizeTimer(ChFinalize)
	go queueing.Enqueue(ChQSHP)
	go queueing.ClearCompletedPeriodically()

	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)
}
