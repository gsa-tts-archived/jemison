ENV?=dev

SHELL=/bin/bash

.PHONY: apply
all: apply

.PHONY: init
init:
	terraform init

.PHONY: clean
clean:
	@echo "------- CLEAN -------"
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	# Don't casually delete the tfstate...
	# rm -f terraform.tfstate
	rm -rf app/
	rm -rf zips/*.zip
	mkdir -p plans

.PHONY: clean_app
clean_app:
	rm -rf app/

.PHONY: app
app: clean_app
	@echo "------- APP -------"
	mkdir -p app
	cd .. ; cp -r config terraform/app/

.PHONY: config
config:
	echo "------- CONFIG -------"
	pushd ../config ; make ; popd

.PHONY: admin
admin: app
	echo "----------------------"
	echo "-- migrate / admin"
	echo "----------------------"
	echo
	# Make the migration tooling
	pushd ../cmd/migrate ; make build ; popd
	cp ../cmd/admin/service.exe app/admin.exe
	cp ../cmd/migrate/service.exe app/migrate.exe
	echo "web: ./migrate.exe && ./admin.exe || exit 1" > app/Procfile
	cp apt.yml app/apt.yml
	pushd app ; zip -r -X ../zips/admin.zip . -x terraform/* ; popd


.PHONY: services
services: app config
	echo "------- SERVICES -------"
	for service in entree extract fetch pack serve walk ; \
	do \
		echo "----------------------" ; \
		echo "-- $${service} " ; \
		echo "----------------------" ; \
		echo ; \
		rm -rf app/ ; \
		mkdir -p app ; \
		pushd .. ; cp -r config terraform/app/ ; popd ; \
		pushd ../cmd/$${service} ; make build ; popd ; \
		cp ../cmd/$${service}/service.exe app/$${service}.exe ; \
		echo "web: ./$${service}.exe || exit 1" > app/Procfile ; \
		cp apt.yml app/apt.yml ; \
		pushd app ; zip -r -X ../zips/$${service}.zip . -x terraform/* ; popd ; \
	done

.PHONY: container_build
container_build:
	echo "------- CONTAINER_BUILD -------"
	docker run -v ${PWD}/..:/app -t jemison/build

# add back container_build
.PHONY: plan
plan: clean init  admin services
	echo "------- PLAN -------"
	cp ${ENV}/variables.tf variables.tf
	terraform plan -var-file="terraform.tfvars" -out plans/das.plan

.PHONY: cfclean
cfclean:
	echo "------- CFCLEAN -------"
	for app in admin entree extract fetch pack serve walk ; do cf delete -f -r $${app} ; done

.PHONY: apply
apply: clean plan apply
	echo "------- APPLY -------"
	terraform apply plans/das.plan

