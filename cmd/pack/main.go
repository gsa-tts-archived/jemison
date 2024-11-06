package main

import (
	"log"
	"net/http"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/sqlite"
)

func main() {
	env.InitGlobalEnv()
	InitializeStorage()
	InitializeQueues()
	engine := common.InitializeAPI()

	log.Println("environment initialized")

	ch_finalize = make(chan *sqlite.PackTable)
	go FinalizeTimer(ch_finalize)

	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)
}
