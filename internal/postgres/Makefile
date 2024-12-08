generate:
	sqlc generate -f search_db-sqlc.yaml
	sqlc generate -f work_db-sqlc.yaml
	
check:
	staticcheck
test:
	go test -tags fts5 *.go -v
