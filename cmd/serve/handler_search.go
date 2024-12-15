package main

import (
	"context"
	"net/http"
	"time"

	"github.com/GSA-TTS/jemison/config"
	"github.com/GSA-TTS/jemison/internal/postgres/search_db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SearchRequestInput struct {
	Host  string `json:"host"`
	Path  string `json:"path"`
	Terms string `json:"terms"`
}

type SearchResult struct {
	Terms      string
	PathID     int64
	PageTitle  string
	Kind       int64
	Weight     float64
	Rank       float64
	Text       string
	PathString string
	Snippet    string
}

func runQuery(sri SearchRequestInput) ([]SearchResult, time.Duration, error) {
	start := time.Now()
	d64, err := config.FQDNToDomain64(sri.Host)
	if err != nil {
		zap.L().Error("could not get Domain64",
			zap.String("host", sri.Host))
		duration := time.Since(start)
		return []SearchResult{}, duration, err
	}

	zap.L().Debug("searching for Domain64 range",
		zap.Int64("start", d64),
		zap.Int64("end", d64+1))

	res, err := JDB.SearchDBQueries.SearchContent(context.Background(),
		search_db.SearchContentParams{
			Query:    sri.Terms,
			D64Start: d64,
			D64End:   d64 + 1,
		})

	duration := time.Since(start)

	cleaned := make([]SearchResult, 0)
	for _, r := range res {
		cleaned = append(cleaned, SearchResult{
			Terms:      sri.Terms,
			PageTitle:  r.Title,
			PathString: r.Path,
			Snippet:    string(r.Snippet),
			Rank:       float64(r.Rank),
		})
	}
	return cleaned, duration, err

}

func SearchHandler(c *gin.Context) {
	var sri SearchRequestInput

	if err := c.BindJSON(&sri); err != nil {
		return
	}

	rows, duration, err := runQuery(sri)

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
