package schemas

import (
	"context"
	"database/sql"
	"strings"

	"go.uber.org/zap"
)

func prepDB(sqlite_filename string) (*sql.DB, error) {

	// FIXME: Any params to the DB?
	db, err := sql.Open("sqlite3", sqlite_filename+"?mode=ro")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	// https://phiresky.github.io/blog/2020/sqlite-performance-tuning/
	db.Exec("pragma journal_mode = WAL")
	db.Exec("pragma synchronous = normal")
	db.Exec("pragma temp_store = file")
	db.Exec("pragma temp_store_directory = /home/vcap/app/tmp")
	// We don't have much RAM. No.
	//db.Exec("pragma mmap_size = 30000000000")
	// Unsure how this effects final filesize or performance on read.
	// db.Exec("pragma page_size = 32768")
	return db, nil
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

type SearchParams struct {
	Terms  []string
	Path   string
	Limit  int64
	Offset int64
}

func NewSearch(terms string) *SearchParams {
	sp := &SearchParams{}
	sp.Terms = strings.Split(terms, " ")
	// By default, match all paths.
	sp.Path = "%%"
	sp.Limit = 10
	sp.Offset = 0
	return sp
}

func (sp *SearchParams) TermsToString() string {
	return strings.Join(sp.Terms[:], " ")
}

// FIXME: There's some optimization possible here --- prepared query,
// etc. This is probably far from optimal.
func (q *Queries) Search(ctx context.Context, params *SearchParams) ([]SearchResult, error) {
	//db, _ := prepDB(file)
	query := `
SELECT 
  path_id,
  (SELECT path from paths WHERE id = path_id), 
  kind,
  weight,
	rank,
	page_title,
	snip
  FROM
    (SELECT 
			titles_fts.path_id as path_id, 
			4.0 as weight, 
			rank, 
			kind, 
			title as txt, 
			title as page_title, 
			snippet(titles_fts, 2, '<b>', '</b>', '...', 8) as snip
      FROM titles_fts
      WHERE titles_fts MATCH ?1
				AND titles_fts.path_id IN (SELECT path_id FROM paths WHERE path LIKE ?2)
    UNION ALL
    SELECT 
			headers_fts.path_id as path_id, 
			2.0 as weight, 
			rank, 
			kind, 
			header as txt, 
			(select title from titles where headers_fts.path_id = titles.path_id) as page_title,
			snippet(headers_fts, 3, '<b>', '</b>', '...', 16) as snip
      FROM headers_fts
      WHERE headers_fts MATCH ?1
				AND headers_fts.path_id IN (SELECT path_id FROM paths WHERE path LIKE ?2)
    UNION ALL
    SELECT 
			bodies_fts.path_id as path_id, 
			1.0 as weight, 
			rank, 
			kind, 
			body as txt, 
			(select title from titles where bodies_fts.path_id = titles.path_id) as page_title,
			snippet(bodies_fts, 2, '<b>', '</b>', '...', 32) as snip
			FROM bodies_fts
      WHERE bodies_fts MATCH ?1
				AND bodies_fts.path_id IN (SELECT path_id FROM paths WHERE path LIKE ?2)
    ORDER BY weight DESC, rank ASC)
		LIMIT ?3
		OFFSET ?4
  ;
`

	zap.L().Debug("running FTS5 query",
		zap.String("params", params.TermsToString()),
		zap.String("path", params.Path),
		zap.Int64("limit", params.Limit),
		zap.Int64("offset", params.Offset))

	rows, err := q.db.QueryContext(ctx,
		query,
		params.TermsToString(),
		params.Path,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		zap.L().Warn("FTS5 search errored",
			zap.String("params", params.TermsToString()),
			zap.String("path", params.Path),
			zap.Int64("limit", params.Limit),
			zap.Int64("offset", params.Offset),
		)
		return nil, err
	}
	defer rows.Close()

	// Always have an empty set to return.
	results := make([]SearchResult, 0)
	// Append any results.
	// The order of the .Scan must match the order of columns
	// defined in the query results above.
	for rows.Next() {
		var sr SearchResult
		sr.Terms = params.TermsToString()
		if err := rows.Scan(
			&sr.PathID,
			&sr.PathString,
			&sr.Kind,
			&sr.Weight,
			&sr.Rank,
			&sr.PageTitle,
			&sr.Snippet,
		); err != nil {
			return results, err
		}
		results = append(results, sr)
	}
	return results, nil
}
