package main

import (
	"fmt"
	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"go.uber.org/zap"
	"net/http"
)

var ThisServiceName = "collect"

func main() {
	fmt.Println("Hello world from collect")
	s, _ := env.Env.GetUserService(ThisServiceName)
	engine := common.InitializeAPI()

	externalHost := s.GetParamString("external_host")
	externalPort := s.GetParamInt64("external_port")

	fmt.Println(ThisServiceName, " environment initialized")

	zap.L().Info("collect environment",
		zap.String("external_host", externalHost),
		zap.Int64("external_port", externalPort),
	)

	zap.L().Info("listening from collect",
		zap.String("port", env.Env.Port))

	// Local and Cloud should both get this from the environment.
	err := http.ListenAndServe(":"+env.Env.Port, engine)
	if err != nil {
		zap.Error(err)
	}
}
