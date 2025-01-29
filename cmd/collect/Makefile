run:
	go run *.go

test:
	go test

check:
	staticcheck

.PHONY: generate
generate:
	echo "nothing to generate"

.PHONY: build
build: generate
	go build -buildvcs=false -o service.exe

clean:
	rm -f service.exe