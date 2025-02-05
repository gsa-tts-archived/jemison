.PHONY: all
all:
	cd domain64 ; make all
	cd services ; make all
	cd hosts ; make all
	mv services/*.yaml .
	mv hosts/*.yaml . 

# .PHONY: system
# system: *.jsonnet
# 	jsonnetfmt \
# 		-i \
# 		--indent 2 \
# 		--max-blank-lines 2 \
# 		--sort-imports \
# 		--string-style s \
# 		--comment-style s \
# 	$<