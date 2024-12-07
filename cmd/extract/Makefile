run:	
	go run *.go

test:
	go test -v -tags fts5 ./...

check:
	staticcheck

generate:
	echo "nothing to generate"
	
build: generate
	go build -buildvcs=false -tags fts5 -o service.exe

clean:
	rm -f fetch.exe
	rm -f Procfile
	rm -f apt.yml
