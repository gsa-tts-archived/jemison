all: clean init plan apply 

init:
	cp ../developers.tf . 
	terraform init

plan:
	mkdir -p plans/
	terraform plan -var-file="terraform.tfvars" -out plans/das.plan

.PHONY: apply
apply:
	echo "------- APPLY -------"
	terraform apply plans/das.plan

.PHONY: clean
clean:
	@echo "------- CLEAN -------"
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -rf plans

.PHONY: cfclean
cfclean:
	echo "------- CFCLEAN -------"
	for app in admin entree extract fetch pack serve walk ; do cf delete -f -r $${app} ; done
	for service in extract fetch jemison-qeueues-db jemison-search-db jemison-work-db queuesdb searchdb workdb serve ; do cf delete-service -f -w $${service} ; done
	