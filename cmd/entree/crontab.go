package main

import (
	"fmt"

	_ "embed"

	"github.com/GSA-TTS/jemison/config"
	"github.com/GSA-TTS/jemison/internal/queueing"

	"github.com/robfig/cron"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

/*
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/

func crontab(schedule string) {
	c := cron.New()
	// https://crontab.guru/#*_*_*_*_*
	err := c.AddFunc("0 * * * * *", section("minutely", schedule))
	if err != nil {
		zap.L().Error("failed to add crontab in entree")
	}
	err = c.AddFunc("@hourly", section("hourly", schedule))
	if err != nil {
		zap.L().Error("failed to add crontab in entree")
	}
	c.Start()
}

func section(section string, schedule string) func() {
	JSON := config.ReadJsonConfig(schedule)
	return func() {
		zap.L().Debug(section)
		for _, site := range gjson.Get(JSON, section).Array() {
			// Clear out any hall passes at this point.
			HallPassLedger.Remove(site.Get("host").String())

			zap.L().Debug(fmt.Sprintln(site))
			ChQSHP <- queueing.QSHP{
				Queue:  "fetch",
				Scheme: site.Get("scheme").String(),
				Host:   site.Get("host").String(),
				Path:   site.Get("path").String(),
			}
		}
	}
}
