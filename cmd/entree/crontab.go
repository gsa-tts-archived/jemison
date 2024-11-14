package main

import (
	"context"
	"fmt"
	"log"

	_ "embed"

	schemas_entree "github.com/GSA-TTS/jemison/cmd/entree/schemas"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/golang-collections/collections/set"
	"github.com/google/go-jsonnet"
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

var JSON = "[]"

func loadJsonnet() {

	vm := jsonnet.MakeVM()
	json, err := vm.EvaluateAnonymousSnippet("schedule.jsonnet", sonnet)
	if err != nil {
		log.Fatal(err)
	}
	JSON = json
}

func crontab() {
	loadJsonnet()

	c := cron.New()
	// https://crontab.guru/#*_*_*_*_*
	c.AddFunc("0 */2 * * * *", loadJsonnet)
	c.AddFunc("0 * * * * *", section("minutely"))
	c.AddFunc("@hourly", section("hourly"))
	c.Start()
}

func section(section string) func() {
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
	hosts := set.New()
	loadJsonnet()
	found := make(map[string]int64)

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
			hosts.Insert(
				site.Get("scheme").String() + "://" +
					site.Get("host").String() +
					site.Get("path").String(),
			)
		}
	}

	hosts.Do(func(host interface{}) {
		id, err := queries.UpsertUniqueHost(ctx, pgtype.Text{String: host.(string), Valid: true})
		if err != nil {
			zap.L().Error("did not get `id` back for host", zap.String("host", host.(string)))
		}

		found[host.(string)] = id
	})

	return found
}
