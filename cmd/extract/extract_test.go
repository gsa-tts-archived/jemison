package main

import (
	"iter"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/GSA-TTS/jemison/internal/sqlite"
	"github.com/GSA-TTS/jemison/internal/sqlite/schemas"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setup(t *testing.T, sqlite_file string) *sqlite.PackTable {
	os.Remove(sqlite_file)
	pt, err := sqlite.CreatePackTable(sqlite_file)
	if err != nil {
		t.Error(err)
	}
	return pt
}

func open(t *testing.T, sqlite_file string) *sqlite.PackTable {
	pt, err := sqlite.OpenPackTable(sqlite_file)
	if err != nil {
		t.Error(err)
	}
	return pt
}

func TestCreateTable(t *testing.T) {
	setup(t, "a.db")
}

func _getLevel(h string) int64 {
	s, _ := strconv.Atoi(h[1:len(h)])
	return int64(s)
}

func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for a := range seq {
			if !yield(f(a)) {
				return
			}
		}
	}
}

func TestExtractHeaders(t *testing.T) {
	pt := setup(t, "headers.db")
	path_id, err := pt.Queries.InsertPath(pt.Context, "/constitution")
	if err != nil {
		t.Error(err)
	}

	fp, err := os.Open("test-files/constitution-02.html")
	if err != nil {
		t.Error(err)
	}
	doc, err := goquery.NewDocumentFromReader(fp)
	if err != nil {
		zap.L().Fatal("cannot create new doc from raw file")
	}

	H := _getHeaders(doc)

	for tag, headers := range H {
		lvl := _getLevel(tag)
		for _, h := range headers {
			id, err := pt.Queries.InsertHeader(pt.Context, schemas.InsertHeaderParams{
				PathID: path_id,
				Level:  lvl,
				Header: h,
			})
			if err != nil {
				t.Error("insert error", err)
			}
			assert.Greater(t, id, int64(0))
		}
	}

	search_params := pt.Queries.NewSearch("north")
	res, _ := pt.Queries.Search(pt.Context, search_params)
	found := false
	for _, r := range res {
		log.Println(r.Text)
		if strings.Contains(r.Text, "north") && strings.Contains(r.Text, "carolina") {
			found = true
		}
	}
	assert.True(t, found)
}
