SHELL=/bin/bash

.PHONY: bundle
bundle: clean generate gz

.PHONY: tar
tar:
	for d in cmd config internal pkg ; do \
		tar --append --file=app.tar \
			../../../$${d} >/dev/null 2>&1 ; \
	done

.PHONY: gz
gz: tar
	gzip app.tar

.PHONY: generate
generate:
	cd ../../../internal/postgres ; make generate

.PHONY: clean
clean:
	cd ../../.. ; make clean
	rm -f app.tar
	rm -f app.tar.gz