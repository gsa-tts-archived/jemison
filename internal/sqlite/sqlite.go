package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"os"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	search_db "github.com/GSA-TTS/jemison/internal/sqlite/schemas"
)

//go:embed schema.sql
var ddl string

// FIXME: This may become an interface?
type PackTable struct {
	Filename     string
	TempFilename string
	DB           *sql.DB
	Queries      *search_db.Queries
}

// Used in testing
func OpenPackTable(db_filename string) (*PackTable, error) {

	pt := PackTable{}
	pt.Filename = db_filename

	// FIXME: Any params to the DB?
	db, err := sql.Open("sqlite3", db_filename+"?_fk=true")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)

	queries := search_db.New(db)

	pt.DB = db
	pt.Queries = queries
	return &pt, nil
}

// FIXME: is it possible for a DB to be open, and be being written to,
// and then the timeout fires again, and we try and write to it again?
// With multiple workers, maybe this is possible, if writing a DB takes
// longer than the timeout period. (Which... when running locally, might be true.)
// So, this needs some protections. Like... maybe a per-host lock, so
// we don't re-open a DB while it is already open.
func CreatePackTable(dbFilename string) (*PackTable, error) {

	ctx := context.Background()
	// Use a temp filename until it is packed.
	tempFilename := uuid.NewString() + ".sqlite"
	//
	os.Remove(tempFilename)

	pt := PackTable{}
	pt.Filename = dbFilename
	pt.TempFilename = tempFilename

	// FIXME: Any params to the DB?
	db, err := sql.Open("sqlite3", tempFilename+"?_fk=true")
	if err != nil {
		zap.L().Error("cannot create SQLite file for packing",
			zap.String("db_filename", tempFilename),
			zap.String("err", err.Error()),
		)
		return nil, err
	}

	db.SetMaxOpenConns(100)
	// https://phiresky.github.io/blog/2020/sqlite-performance-tuning/
	_, err = db.Exec("pragma journal_mode = WAL")
	if err != nil {
		zap.L().Error("pragma fail journal_mode", zap.String("err", err.Error()))
	}
	_, err = db.Exec("pragma synchronous = normal")
	if err != nil {
		zap.L().Error("pragma fail synchronous", zap.String("err", err.Error()))
	}
	_, err = db.Exec("pragma temp_store = file")
	if err != nil {
		zap.L().Error("pragma fail temp_store", zap.String("err", err.Error()))
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := search_db.New(db)

	pt.DB = db
	pt.Queries = queries

	return &pt, nil
}

func (pt *PackTable) PrepForNetwork() {
	// https://turso.tech/blog/something-you-probably-want-to-know-about-if-youre-using-sqlite-in-golang-72547ad625f1
	db, _ := sql.Open("sqlite3", pt.TempFilename)
	pt.DB = db
	_, err := pt.DB.ExecContext(context.Background(), "PRAGMA wal_checkpoint(TRUNCATE)")
	if err != nil {
		zap.L().Error("pragma fail on prep for network truncate")
	}
	_, err = pt.DB.ExecContext(context.Background(), "PRAGMA vacuum")
	if err != nil {
		zap.L().Error("pragma fail on prep for network vacuum")
	}
	_, err = pt.DB.ExecContext(context.Background(), "PRAGMA optimize")
	if err != nil {
		zap.L().Error("pragma fail on prep for network optimize")
	}

	err = pt.DB.Close()

	// Now, rename the temp to the real thing.
	os.Rename(pt.TempFilename, pt.Filename)

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
