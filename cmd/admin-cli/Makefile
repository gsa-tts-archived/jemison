clean:
	rm -f service.exe
	rm -f __debug*
	
run: clean
	jsonnet config.jsonnet > config.json
	go run *.go