package sqlite

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"testing"

	"github.com/GSA-TTS/jemison/internal/sqlite/schemas"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tjarratt/babble"
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
		Title:  strings.ToLower("The Most Incredible wind windy winding Title"),
	})
}

func TestInsertBadTitle(t *testing.T) {
	pt := open(t, "many_paths.db")
	_, err := pt.Queries.InsertTitle(pt.Context, schemas.InsertTitleParams{
		PathID: 0,
		Title:  strings.ToLower("The Most Incredible Title"),
	})
	if err != nil {
		assert.True(t, strings.Contains(err.Error(), "FOREIGN KEY"))
	}
}

func TestInsertManyTitles(t *testing.T) {
	pt := open(t, "many_paths.db")
	babbler := babble.NewBabbler()
	babbler.Count = 3
	babbler.Separator = " "
	for ndx := range 100 {
		id, err := pt.Queries.InsertTitle(pt.Context, schemas.InsertTitleParams{
			PathID: int64(ndx + 1),
			Title:  strings.ToLower(fmt.Sprintf("Title Number %d: "+babbler.Babble(), ndx+2)),
		})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, ndx+2, id)
	}
}

func TestInsertManyHeaders(t *testing.T) {
	pt := open(t, "many_paths.db")
	babbler := babble.NewBabbler()
	babbler.Count = 3
	babbler.Separator = " "
	for ndx := range 100 {
		windy := ""
		if rand.IntN(100) < 5 {
			windy = "windy"
		}

		for range rand.IntN(6) {
			_, err := pt.Queries.InsertHeader(pt.Context, schemas.InsertHeaderParams{
				PathID: int64(ndx + 1),
				Level:  int64(rand.IntN(6) + 1),
				Header: fmt.Sprintf("%s Title Number %d: "+babbler.Babble(), windy, ndx+2),
			})
			if err != nil {
				t.Error(err)
			}
			//assert.Equal(t, ndx+2, id+int64(cnt))
		}
	}
}

func TestInsertContent(t *testing.T) {
	babbler := babble.NewBabbler()
	babbler.Count = 50
	babbler.Separator = " "
	pt := open(t, "many_paths.db")
	for ndx := range 100 {
		id, err := pt.Queries.InsertBody(pt.Context, schemas.InsertBodyParams{
			PathID: int64(rand.IntN(90) + 1),
			Body:   strings.ToLower(babbler.Babble()),
			Tag:    "p",
		})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, ndx+1, id)
	}
}

func TestSearch(t *testing.T) {
	pt := open(t, "many_paths.db")
	r, err := pt.Queries.Search(pt.Context, schemas.NewSearch("wind*"))
	if err != nil {
		t.Error(err)
	}
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"#", "PathID", "Kind", "Weight", "Rank"})
	for ndx, res := range r {
		tab.AppendRow([]interface{}{ndx, res.PathID, res.Kind, res.Weight, res.Rank})
	}
	tab.Render()
	// We should always get back the zeroth title, because
	// the word `wind` is in it. And `windy`. And `winding`.
	assert.Equal(t, 1, r[0].PathID)
}
