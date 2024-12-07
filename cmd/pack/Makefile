run:	
	go run --tags fts5 *.go

test:
	go test

check:
	staticcheck

generate:
	echo "nothing to generate"
	
build: generate
	go build --tags fts5 -buildvcs=false -o service.exe

clean:
	rm -f service.exe
