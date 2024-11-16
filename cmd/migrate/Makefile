run:	
	go run *.go

test:
	go test

check:
	staticcheck

# -tags=viper_bind_struct 
build: 
	go build -buildvcs=false -o service.exe

clean:
	rm -f service.exe
	rm -f schemas/*.go
