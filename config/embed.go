package config

import (
	"embed"
	_ "embed"
	"log"

	"github.com/google/go-jsonnet"
)

//go:embed *.jsonnet
var ConfigFs embed.FS

func ReadConfigJsonnet(sonnetFilename string) string {
	bytes, _ := ConfigFs.ReadFile(sonnetFilename)

	vm := jsonnet.MakeVM()
	json, err := vm.EvaluateAnonymousSnippet(sonnetFilename, string(bytes))
	if err != nil {
		log.Fatal(err)
	}
	return json
}
