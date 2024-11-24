package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	search_db "github.com/GSA-TTS/jemison/internal/sqlite/schemas"
)

//go:embed schema.sql
var ddl string

// FIXME: This may become an interface?
type PackTable struct {
	Filename string
	Context  context.Context
	DB       *sql.DB
	Queries  *search_db.Queries
}

func OpenPackTable(db_filename string) (*PackTable, error) {

	pt := PackTable{}
	pt.Filename = db_filename

	ctx := context.Background()

	// FIXME: Any params to the DB?
	db, err := sql.Open("sqlite3", db_filename+"?_fk=true")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)

	queries := search_db.New(db)

	pt.Context = ctx
	pt.DB = db
	pt.Queries = queries
	return &pt, nil
}

func CreatePackTable(db_filename string) (*PackTable, error) {

	pt := PackTable{}
	pt.Filename = db_filename

	ctx := context.Background()
	// What if the file already exists? We should clean up first.
	// When testing locally, it is the same location as where we serve from.
	// In cloud.gov, this will be on a pack instance.
	os.Remove(db_filename)

	// FIXME: Any params to the DB?
	db, err := sql.Open("sqlite3", db_filename+"?_fk=true")
	if err != nil {
		zap.L().Error("cannot create SQLite file for packing",
			zap.String("db_filename", db_filename),
			zap.String("err", err.Error()),
		)
		return nil, err
	}
	db.SetMaxOpenConns(100)
	// https://phiresky.github.io/blog/2020/sqlite-performance-tuning/
	_, err = db.Exec("pragma journal_mode = WAL")
	if err != nil {
		zap.L().Error("pragma fail on sqlite")
	}
	_, err = db.Exec("pragma synchronous = normal")
	if err != nil {
		zap.L().Error("pragma fail on sqlite")
	}
	_, err = db.Exec("pragma temp_store = file")
	if err != nil {
		zap.L().Error("pragma fail on sqlite")
	}
	_, err = db.Exec("pragma temp_store_directory = /home/vcap/app/tmp")
	if err != nil {
		zap.L().Error("pragma fail on sqlite")
	}

	// We don't have much RAM. No.
	//db.Exec("pragma mmap_size = 30000000000")
	// Unsure how this effects final filesize or performance on read.
	// db.Exec("pragma page_size = 32768")
	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := search_db.New(db)

	pt.Context = ctx
	pt.DB = db
	pt.Queries = queries

	return &pt, nil
}

func (pt *PackTable) PrepForNetwork() {
	// https://turso.tech/blog/something-you-probably-want-to-know-about-if-youre-using-sqlite-in-golang-72547ad625f1
	db, _ := sql.Open("sqlite3", pt.Filename)
	pt.DB = db
	_, err := pt.DB.ExecContext(pt.Context, "PRAGMA wal_checkpoint(TRUNCATE)")
	if err != nil {
		zap.L().Error("pragma fail on prep for network truncate")
	}
	_, err = pt.DB.ExecContext(pt.Context, "PRAGMA vacuum")
	if err != nil {
		zap.L().Error("pragma fail on prep for network vacuum")
	}
	_, err = pt.DB.ExecContext(pt.Context, "PRAGMA optimize")
	if err != nil {
		zap.L().Error("pragma fail on prep for network optimize")
	}

	err = pt.DB.Close()
	if err != nil {
		zap.L().Error("pragma fail on sqlite close")
	}

}

func SqliteFilename(db_filename string) string {
	// Always add an .sqlite extension to filenames.
	if has_ext := strings.HasSuffix(db_filename, "sqlite"); has_ext {
		//zap.L().Debug("not adding .sqlite to filename")
		return db_filename
	} else {
		//zap.L().Debug("adding .sqlite to filename")
		db_filename = db_filename + ".sqlite"
		return db_filename
	}
}
