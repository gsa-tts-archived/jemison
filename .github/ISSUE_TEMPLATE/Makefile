fmt: $(wildcard *.jsonnet)
	jsonnetfmt \
		-i \
		--indent 2 \
		--max-blank-lines 2 \
		--sort-imports \
		--string-style s \
		--comment-style s \
		$?

build: $(wildcard *.jsonnet)
	# $(subst from,to,text)
	jsonnet $? > $(subst jsonnet,yaml,$?)

all: fmt build
.PHONY: all