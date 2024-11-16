package sqlite

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/GSA-TTS/jemison/internal/sqlite/schemas"
	"github.com/zeebo/assert"
)

func setup(t *testing.T, sqlite_file string) *PackTable {
	os.Remove(sqlite_file)
	pt, err := CreatePackTable(sqlite_file)
	if err != nil {
		t.Error(err)
	}
	return pt
}

func open(t *testing.T, sqlite_file string) *PackTable {
	pt, err := OpenPackTable(sqlite_file)
	if err != nil {
		t.Error(err)
	}
	return pt
}

func TestCreateTable(t *testing.T) {
	setup(t, "a.db")
}

func TestInsertPath(t *testing.T) {
	pt := setup(t, "one_path.db")
	id, _ := pt.Queries.InsertPath(pt.Context, "/a")
	assert.Equal(t, 1, id)
	pt.PrepForNetwork()
}

// This sequence of tests builds up one database.
// This is to test FK constraints, amongst other things.
func TestInsertManyPaths(t *testing.T) {
	pt := setup(t, "many_paths.db")
	for i := range 100 {
		id, _ := pt.Queries.InsertPath(pt.Context, fmt.Sprintf("/path/%d", i+1))
		assert.Equal(t, i+1, id)
	}
	pt.PrepForNetwork()
}

func TestInsertOneTitle(t *testing.T) {
	pt := open(t, "many_paths.db")
	pt.Queries.InsertTitle(pt.Context, schemas.InsertTitleParams{
		PathID: 1,
		Txt:    "The Most Incredible Title",
	})
}

func TestInsertBadTitle(t *testing.T) {
	pt := open(t, "many_paths.db")
	_, err := pt.Queries.InsertTitle(pt.Context, schemas.InsertTitleParams{
		PathID: 0,
		Txt:    "The Most Incredible Title",
	})
	if err != nil {
		assert.True(t, strings.Contains(err.Error(), "FOREIGN KEY"))
	}
}

func TestInsertManyTitles(t *testing.T) {
	pt := open(t, "many_paths.db")
	for ndx := range 100 {
		id, err := pt.Queries.InsertTitle(pt.Context, schemas.InsertTitleParams{
			PathID: int64(ndx + 1),
			Txt:    fmt.Sprintf("Title Number %d", ndx+2),
		})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, ndx+2, id)
	}
}
