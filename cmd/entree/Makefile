# -tags=viper_bind_struct 
build:
	go build -buildvcs=false -o service.exe

test:
	go test

check:
	staticcheck

generate:
	sqlc generate || exit 1


clean:
	rm -f service.exe
	rm -f schemas/*.go
