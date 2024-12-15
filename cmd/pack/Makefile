# --tags fts5 
run:	
	go run *.go

test:
	go test

check:
	staticcheck

generate:
	echo "nothing to generate"
	
# --tags fts5 
build: generate
	go build -buildvcs=false -o service.exe

clean:
	rm -f service.exe
