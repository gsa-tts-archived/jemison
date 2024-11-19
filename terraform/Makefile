.PHONY: init
init:
	terraform init

.PHONY: clean
clean:
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -rf app/
	rm -rf zips/*.zip
	mkdir -p plans

.PHONY: clean_app
clean_app:
	rm -rf app/

.PHONY: app
app: clean_app
	mkdir -p app
	cd .. ; \
		cp -r assets terraform/app/ ; \
		cp -r config terraform/app/

.PHONY: admin
admin: app
	mkdir -p zips
	cp ../cmd/admin/service.exe app/admin.exe
	echo "web: ./admin.exe" > app/Procfile
	cp ../assets/apt.yml app/apt.yml
	cd app ; zip -r -X ../zips/admin.zip . \
		-x assets/databases/* \
		-x assets/static.zip \
		-x assetts/static/* \
		-x terraform/*



.PHONY: entree
entree: app
	mkdir -p zips
	cp ../cmd/entree/service.exe app/extract.exe
	cp ../cmd/entree/run.sh app/run.sh
	echo "web: ./run.sh" > app/Procfile
	cp ../assets/apt.yml app/apt.yml
	cp ../cmd/migrate/service.exe app/migrate.exe
	cd app ; zip -r -X ../zips/entree.zip . \
		-x assets/databases/* \
		-x assets/static.zip \
		-x assetts/static/* \
		-x terraform/*

.PHONY: extract
extract: app
	mkdir -p zips
	cp ../cmd/extract/service.exe app/extract.exe
	echo "web: ./extract.exe" > app/Procfile
	cp ../assets/apt.yml app/apt.yml
	cd app ; zip -r -X ../zips/extract.zip . \
		-x assets/databases/* \
		-x assets/static.zip \
		-x assetts/static/* \
		-x terraform/*

.PHONY: fetch
fetch: app
	mkdir -p zips
	cp ../cmd/fetch/service.exe app/fetch.exe
	echo "web: ./fetch.exe" > app/Procfile
	cp ../assets/apt.yml app/apt.yml
	cd app ; zip -r -X ../zips/fetch.zip . \
		-x assets/databases/* \
		-x assets/static.zip \
		-x assetts/static/* \
		-x terraform/*

# .PHONY: migrate
# migrate: app
# 	cp ../cmd/migrate/service.exe app/migrate.exe
# 	echo "web: ./migrate.exe" > app/Procfile
# 	cp ../assets/apt.yml app/apt.yml
# 	cd app ; zip -r -X ../zips/migrate.zip . \
# 		-x assets/databases/* \
# 		-x assets/static.zip \
# 		-x assetts/static/* \
# 		-x terraform/*


.PHONY: pack
pack: app
	mkdir -p zips
	cp ../cmd/pack/service.exe app/pack.exe
	echo "web: ./pack.exe" > app/Procfile
	cp ../assets/apt.yml app/apt.yml
	cd app ; zip -r -X ../zips/pack.zip . \
		-x assets/databases/* \
		-x assets/static.zip \
		-x assets/static/* \
		-x terraform/*

.PHONY: serve
serve: app
	mkdir -p zips
	cp ../cmd/serve/service.exe app/serve.exe
	echo "web: ./serve.exe" > app/Procfile
	cp ../assets/apt.yml app/apt.yml
	cp ../assets/static.zip app/static.zip
	cd app ; unzip -qq static.zip ; rm static.zip 
	mkdir app/databases ;
	cp ../assets/static/index.html app/static/index.html
	cd app ; zip -r -X ../zips/serve.zip . \
		-x assets/databases/*sqlite* \
		-x assets/static.zip \
		-x terraform/*

.PHONY: walk
walk: app
	mkdir -p zips
	cp ../cmd/walk/service.exe app/walk.exe
	echo "web: ./walk.exe" > app/Procfile
	cp ../assets/apt.yml app/apt.yml
	cd app ; zip -r -X ../zips/walk.zip . \
		-x assets/databases/* \
		-x assets/static.zip \
		-x assetts/static/* \
		-x terraform/*

.PHONY:
plan:
	terraform plan -var-file="terraform.tfvars" -out plans/das.plan

.PHONY: tf_admin
tf_admin: admin plan apply
	echo "tf admin"

.PHONY: tf_entree
tf_entree: entree plan apply
	echo "tf entree"

.PHONY: tf_migrate
tf_migrate: migrate plan apply
	echo "tf migrate"

.PHONY: tf_extract
tf_extract: extract plan apply
	echo "tf extract"

.PHONY: tf_fetch
tf_fetch: plan apply
	echo "tf fetch"

.PHONY: tf_pack
tf_pack: pack plan apply
	echo "tf pack"

.PHONY: tf_serve
tf_serve: serve plan apply
	echo "tf serve"

.PHONY: tf_walk
tf_walk: walk plan apply
	echo "tf walk"
	
.PHONY: container_build
container_build:
	docker run -v ${PWD}:/app -t jemison/build

.PHONY: plan_all
plan_all: init admin entree extract fetch pack serve walk
	terraform plan -var-file="terraform.tfvars" -out plans/das.plan

.PHONY: apply
apply:
	terraform apply plans/das.plan

.PHONY: apply_all
apply_all: clean plan_all apply
	echo "apply_all"

