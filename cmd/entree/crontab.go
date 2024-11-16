package main

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/GSA-TTS/jemison/config"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/work_db/work_db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

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

func crontab() {
	c := cron.New()
	// https://crontab.guru/#*_*_*_*_*
	c.AddFunc("0 * * * * *", section("minutely"))
	c.AddFunc("@hourly", section("hourly"))
	c.Start()
}

func section(section string) func() {
	JSON := config.ReadConfigJsonnet("schedule.jsonnet")
	return func() {
		zap.L().Debug(section)
		for _, site := range gjson.Get(JSON, section).Array() {
			// Clear out any hall passes at this point.
			HallPassLedger.Remove(site.Get("host").String())

			zap.L().Debug(fmt.Sprintln(site))
			queueing.InsertFetch(
				site.Get("scheme").String(),
				site.Get("host").String(),
				site.Get("path").String(),
			)
		}
	}
}

func upsertUniqueHosts() map[string]int64 {
	JSON := config.ReadConfigJsonnet("schedule.jsonnet")
	uniqueHosts := make(map[string]int64)

	ctx := context.Background()

	db_string, err := env.Env.GetDatabaseUrl(env.JemisonWorkDatabase)
	if err != nil {
		zap.L().Fatal("could not get db URL for DB1")
	}
	conn, err := pgx.Connect(ctx, db_string)
	if err != nil {
		zap.L().Fatal("could not connect to DB1")
	}
	defer conn.Close(ctx)

	queries := work_db.New(conn)

	for _, section := range gjson.Parse(JSON).Get("@keys").Array() {
		for _, site := range gjson.Get(JSON, section.String()).Array() {
			// We should never see a -1 in the host table. Not sure
			// how else to do this. The following loop will either populate
			// the value or fail.
			uniqueHosts[site.Get("host").String()] = -1
		}
	}

	// Iterate through the set, and create a unique map of ... oh. I could have
	// just created the map in the first place... FIXME later...
	for h := range uniqueHosts {
		id, err := queries.UpsertUniqueHost(ctx, pgtype.Text{String: h, Valid: true})
		if err != nil {
			zap.L().Error("did not get `id` back for host",
				zap.String("host", h),
				zap.String("err", err.Error()))
		}
		uniqueHosts[h] = id
	}

	return uniqueHosts
}
