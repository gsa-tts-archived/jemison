package util

import (
	"log"

	"github.com/google/go-jsonnet"
)

func ProcessJsonnet(sonnet string) string {

	vm := jsonnet.MakeVM()
	json, err := vm.EvaluateAnonymousSnippet("schedule.jsonnet", sonnet)
	if err != nil {
		log.Fatal(err)
	}
	return json
}
