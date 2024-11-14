package main

import (
	"embed"
	"fmt"
	"net/url"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"go.uber.org/zap"
)

// Carry our migrations with us as part of the build.
// This eliminates wondering where they are when we deploy.
//
//go:embed db/migrations/*.sql
var fs embed.FS

// Assumes config has been read
func Migrate() {

	db1_url, err := env.Env.GetDatabaseUrl(env.DB1)
	if err != nil {
		zap.L().Fatal("could not get url for DB1")
	}

	u, _ := url.Parse(db1_url)
	db := dbmate.New(u)
	db.FS = fs

	fmt.Println("Migrations:")
	migrations, err := db.FindMigrations()
	if err != nil {
		panic(err)
	}
	for _, m := range migrations {
		fmt.Println(m.Version, m.FilePath)
	}

	fmt.Println("\nApplying...")
	err = db.CreateAndMigrate()
	if err != nil {
		panic(err)
	}
}
