run:	
	go run *.go

test:
	go test

check:
	staticcheck

copy:
	cp -R ../../sql/search_db . 
	cp -R ../../sql/work_db . 

# -tags=viper_bind_struct 
build: clean copy
	go build -buildvcs=false -o service.exe

clean:
	rm -f service.exe
	rm -rf search_db
	rm -rf work_db
