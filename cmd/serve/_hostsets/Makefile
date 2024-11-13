SOURCES = $(wildcard gov/*.jsonnet)
POLITE_SLEEP ?= 2

all: clean hostset $(patsubst %.jsonnet,%.json, ${SOURCES})

hostset:
	jsonnet hostsets.jsonnet > hostsets.json

clean: 
	rm -f *.yaml

%.json: %.jsonnet
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
