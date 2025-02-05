SOURCES = $(wildcard *.jsonnet)

all: clean $(patsubst %.jsonnet,%.yaml, ${SOURCES})

clean: 
	rm -f *.json

%.yaml: %.jsonnet
	@echo "IN" $<  "OUT" $@
	jsonnetfmt \
		-i \
		--indent 2 \
		--max-blank-lines 2 \
		--sort-imports \
		--string-style s \
		--comment-style s \
		$<
	jsonnet $< > $@
