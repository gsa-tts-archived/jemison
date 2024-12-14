package main

import (
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"go.uber.org/zap"
)

var ThisServiceName = "migrate"

func main() {
	env.InitGlobalEnv(ThisServiceName)
	MigrateJemisonDB()
	queueing.RunRiverMigrator()

	zap.L().Info("migrations completed successfully")
}
