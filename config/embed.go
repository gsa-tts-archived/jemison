package config

import (
	"embed"

	"github.com/google/go-jsonnet"
	"go.uber.org/zap"
)

//go:embed *.jsonnet
//go:embed *.json
var ConfigFs embed.FS

func ReadConfigJsonnet(sonnetFilename string) string {
	bytes, _ := ConfigFs.ReadFile(sonnetFilename)

	vm := jsonnet.MakeVM()
	json, err := vm.EvaluateAnonymousSnippet(sonnetFilename, string(bytes))
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return json
}

func ReadJsonConfig(jsonFilename string) string {
	json_bytes, err := ConfigFs.ReadFile(jsonFilename)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return string(json_bytes)
}
