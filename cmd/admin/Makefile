run:	
	go run *.go

test:
	go test

check:
	staticcheck

generate:
	cd ../../internal/sqlite ; sqlc generate || exit 1
	cd ../../internal/postgres ; sqlc generate || exit 1

build: clean generate
	go build -buildvcs=false -o service.exe

clean:
	rm -f service.exe
