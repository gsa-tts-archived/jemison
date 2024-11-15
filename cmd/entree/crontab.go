package main

import (
	"context"
	"fmt"
	"net/url"

	_ "embed"

	schemas_entree "github.com/GSA-TTS/jemison/cmd/entree/schemas"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/robfig/cron"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//go:embed schedule.jsonnet
var sonnet string

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
	JSON := util.ProcessJsonnet(sonnet)
	return func() {
		zap.L().Debug(section)
		for _, site := range gjson.Get(JSON, section).Array() {
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
	JSON := util.ProcessJsonnet(sonnet)
	uniqueHosts := make(map[string]int64)

	ctx := context.Background()

	db_string, err := env.Env.GetDatabaseUrl(env.DB1)
	if err != nil {
		zap.L().Fatal("could not get db URL for DB1")
	}
	conn, err := pgx.Connect(ctx, db_string)
	if err != nil {
		zap.L().Fatal("could not connect to DB1")
	}
	defer conn.Close(ctx)

	queries := schemas_entree.New(conn)

	for _, section := range gjson.Parse(JSON).Get("@keys").Array() {
		for _, site := range gjson.Get(JSON, section.String()).Array() {
			u := url.URL{
				Scheme: site.Get("scheme").String(),
				Host:   site.Get("host").String(),
				Path:   site.Get("path").String(),
			}
			uniqueHosts[u.Host] = -1
		}
	}

	// Iterate through the set, and create a unique map of ... oh. I could have
	// just created the map in the first place... FIXME later...
	for h := range uniqueHosts {
		id, err := queries.UpsertUniqueHost(ctx, pgtype.Text{String: h, Valid: true})
		if err != nil {
			zap.L().Error("did not get `id` back for host", zap.String("host", h))
		}
		uniqueHosts[h] = id
	}

	return uniqueHosts
}
