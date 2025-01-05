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
