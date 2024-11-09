SOURCES = $(wildcard *.jsonnet)
DESTINATIONS = $(patsubst %.jsonnet, %.yaml, $(SOURCES))

clean: 
	rm -f *.yaml

jsonnet: $(SOURCES)
	@echo $?
	jsonnetfmt \
		-i \
		--indent 2 \
		--max-blank-lines 2 \
		--sort-imports \
		--string-style s \
		--comment-style s \
		$?
	jsonnet $? > $?.yaml

$(EXEC): $(DESTINATIONS)
		echo $^