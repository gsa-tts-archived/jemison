package main

import (
	"github.com/GSA-TTS/jemison/internal/env"
	"go.uber.org/zap"
)

var ThisServiceName = "migrate"

func main() {
	env.InitGlobalEnv(ThisServiceName)
	MigrateJemisonDB()
	zap.L().Info("migrations completed successfully")
}
