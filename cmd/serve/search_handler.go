package main

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/sqlite/schemas"
	"github.com/gin-gonic/gin"
	"github.com/kljensen/snowball"
	"go.uber.org/zap"
)

// FIXME This becomes the API search interface
type ServeRequestInput struct {
	Host  string `json:"host"`
	Path  string `json:"path"`
	Terms string `json:"terms"`
}

var statmap sync.Map

func runQuery(sri ServeRequestInput) (
	//[]schemas.SearchSiteIndexSnippetsRow,
	[]schemas.SearchResult,
	time.Duration,
	error) {
	start := time.Now()

	s, _ := env.Env.GetUserService("serve")
	database_files_path := s.GetParamString("database_files_path")
	results_per_query := s.GetParamInt64("results_per_query")
	destination := database_files_path + "/" + sri.Host + ".sqlite"

	// Don't only use the stemmed words
	existing_terms := strings.Split(sri.Terms, " ")
	improved_terms := make([]string, 0)
	for _, t := range existing_terms {
		// start by adding the existing terms to the list.
		improved_terms = append(improved_terms, t)
		stemmed, err := snowball.Stem(t, "english", true)
		if err != nil {
			// Pass. Keep the value as-is
			improved_terms = append(improved_terms, t)
		} else {
			if len(stemmed) > 0 {
				// Stem the terms and add wildcards.
				improved_terms = append(improved_terms, stemmed+"*")
			}
		}
	}

	// Definitely need to OR the terms here.
	// Otherwise, we're building a large implicit AND
	improved_terms_string := strings.Join(improved_terms, " OR ")

	zap.L().Debug("search string",
		zap.String("original", sri.Terms),
		zap.String("improved", improved_terms_string))

	db, err := sql.Open("sqlite3", destination+"?cache=shared&mode=ro")
	if err != nil {
		zap.L().Fatal("could not open db connection", zap.String("database", destination))
	}
	db.SetMaxOpenConns(200)
	db.SetConnMaxIdleTime(5000 * time.Millisecond)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(10000 * time.Millisecond)

	path := sri.Path + "%"
	queries := schemas.New(db)
	search_params := schemas.NewSearch(improved_terms_string)
	search_params.Limit = results_per_query
	search_params.Path = path
	// FIXME: This wants to return snippets
	res, err := queries.Search(context.Background(), search_params)
	db.Close()

	// This is all fine and good, but it would be nice to annotate
	// each search result with the terms that were used
	//zap.L().Info("search results", zap.Int("count", len(res)))

	// if (len(res) < 2) && (len(improved_terms) > 1) {
	// 	res = permuteSubqueries(queries, path, improved_terms, results_per_query)
	// }

	duration := time.Since(start)
	//return res[0:min(limit, len(res))], duration, err
	// THIS QUIETS THE LINTER FOR A MOMENT...

	return res, duration, err
}

// //////////////////////////
// Search Handler
// Handles the API requests as they come in.

func SearchHandler(c *gin.Context) {
	var sri ServeRequestInput

	if err := c.BindJSON(&sri); err != nil {
		return
	}

	rows, duration, err := runQuery(sri)
	runStats(sri, duration)

	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"result":  "err",
			"message": err.Error(),
			"elapsed": duration,
			"results": nil,
		})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{
			"result":  "ok",
			"elapsed": duration,
			"results": rows,
		})
		return
	}
}
