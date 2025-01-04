init:
	cp ../developers.tf . 
	terraform init

plan:
	mkdir -p plans/
	terraform plan -var-file="terraform.tfvars" -out plans/das.plan

apply:
	echo "APPLY"

.PHONY: clean
clean:
	@echo "------- CLEAN -------"
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -rf plans