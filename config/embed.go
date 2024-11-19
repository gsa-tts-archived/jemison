package config

import (
	"bytes"
	"embed"

	"github.com/google/go-jsonnet"
	"go.uber.org/zap"
)

//go:embed *.jsonnet
//go:embed *.json
//go:embed *.yaml
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

func GetYamlFileReader(yamlFilename string) *bytes.Reader {
	yaml_bytes, err := ConfigFs.ReadFile(yamlFilename)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return bytes.NewReader(yaml_bytes)
}
