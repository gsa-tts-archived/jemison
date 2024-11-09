SOURCES = $(wildcard *.jsonnet)

clean: 
	rm -f *.yaml

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

all: $(patsubst %.jsonnet,%.yaml, ${SOURCES})