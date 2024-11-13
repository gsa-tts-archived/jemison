.PHONY: run
run: clean
	jsonnet config.jsonnet > config.json
	go run *.go

.PHONY: clean
clean:
	rm -f service.exe
	rm -f __debug*
	