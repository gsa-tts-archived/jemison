package main

import (
	"github.com/GSA-TTS/jemison/internal/env"
	"go.uber.org/zap"
)

var ThisServiceName = "migrate"
var HostIdMap = make(map[string]int64)

func main() {
	env.InitGlobalEnv(ThisServiceName)
	MigrateJemisonDB()
	zap.L().Info("migrations completed successfully")
}
