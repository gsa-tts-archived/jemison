package schemas

import (
	"context"
	"database/sql"

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
	PathID     int64
	PathString string
	Kind       int64
	Weight     float64
	Rank       float64
	Text       string
}

// FIXME: There's some optimization possible here --- prepared query,
// etc. This is probably far from optimal.
func (q *Queries) Search(ctx context.Context, terms string) ([]SearchResult, error) {
	//db, _ := prepDB(file)
	query := `
SELECT 
  path_id,
  (SELECT path from paths WHERE id = path_id), 
  kind,
  weight,
	rank,
  txt
  FROM
    (SELECT titles_fts.path_id as path_id, 4.0 as weight, rank, kind, title as txt
      FROM titles_fts
      WHERE title MATCH ?1
    UNION ALL
    SELECT headers_fts.path_id as path_id, 2.0 as weight, rank, kind, header as txt
      FROM headers_fts
      WHERE header MATCH ?1
    UNION ALL
    SELECT bodies_fts.path_id as path_id, 1.0 as weight, rank, kind, body as txt
      FROM bodies_fts
      WHERE body MATCH ?1
    ORDER BY weight DESC, rank ASC)
  ;
`
	rows, err := q.db.QueryContext(ctx, query, terms)
	if err != nil {
		zap.L().Error("FTS5 search errored", zap.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult

	for rows.Next() {
		var sr SearchResult
		if err := rows.Scan(
			&sr.PathID,
			&sr.PathString,
			&sr.Kind,
			&sr.Weight,
			&sr.Rank,
			&sr.Text); err != nil {
			return results, err
		}
		results = append(results, sr)
	}
	return results, nil
}
