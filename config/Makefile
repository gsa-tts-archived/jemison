clean: $(wildcard .yaml)
	rm -f $?

jsonnet: container.jsonnet, cf.jsonnet
	# $(subst from,to,text)
	jsonnet $? > $(subst "jsonnet","yaml",$?)

all: clean jsonnet
.PHONY: all