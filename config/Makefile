.PHONY: all
all: system
	cd services ; make all
	mv services/*.yaml . 

.PHONY: system
system: *.jsonnet
	jsonnetfmt \
		-i \
		--indent 2 \
		--max-blank-lines 2 \
		--sort-imports \
		--string-style s \
		--comment-style s \
	$<