run:	
	go run *.go

test:
	go test

check:
	staticcheck

generate:
	cd ../../internal/postgres ; make generate || exit 1
	
build: clean generate
	go build -buildvcs=false -o service.exe

clean:
	rm -f service.exe
