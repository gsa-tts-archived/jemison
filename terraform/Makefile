SHELL=/bin/bash

.PHONY: init
init:
	terraform init

.PHONY: clean
clean:
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
	mkdir -p app
	cd .. ; cp -r config terraform/app/

.PHONY: config
config:
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

.PHONY:
plan: init services
	terraform plan -var-file="terraform.tfvars" -out plans/das.plan

.PHONY: container_build
container_build:
	docker run -v ${PWD}:/app -t jemison/build

.PHONY: plan_all
plan_all: init admin services
	terraform plan -var-file="terraform.tfvars" -out plans/das.plan

.PHONY: apply
apply:
	terraform apply plans/das.plan

.PHONY: cfclean
cfclean:
	for app in admin entree extract fetch pack serve walk ; do cf delete -f -r $${app} ; done

.PHONY: apply_all
apply_all: clean cfclean  plan_all apply
	echo "apply_all"

