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
	# cd hostsets ; make

clean:
	rm -f fetch.exe
	rm -f Procfile
	rm -f apt.yml

terraform: clean build
	mkdir -p app
	echo "web: ./service.exe" > Procfile
	zip -r -X app.zip service.exe Procfile apt.yml
