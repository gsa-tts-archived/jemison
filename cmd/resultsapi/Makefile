run:	
	go run *.go

test:
	go test

.PHONY: generate
generate:
	echo "nothing to generate"
	
.PHONY: build
build: generate
	go build -buildvcs=false -o service.exe

clean:
	rm -f search.exe