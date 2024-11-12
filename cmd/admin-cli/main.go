package main

import (
	"log"
	"os"
	"time"

	"github.com/GSA-TTS/jemison/cmd/admin-cli/pages"
	"github.com/spf13/viper"
)

const refreshInterval = 1000 * time.Millisecond

var ThisServiceName = "admin_cli"

func main() {
	os.Setenv("ENV", "LOCALHOST")

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("cmd/admin-cli")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Could not read config.")
	}

	g := pages.NewGui()
	g.MainPage("main")
	g.CrawlPage("crawl")
	g.MonitorJobQueues("monitor")
	g.Run()
}
