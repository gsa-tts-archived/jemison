.PHONY: clean
clean:
	rm -f internal/sqlite/schemas/db.go
	rm -f internal/sqlite/schemas/models.go
	rm -f internal/sqlite/schemas/query.sql.go
	rm -f internal/work_db/work_db/db.go
	rm -f internal/work_db/work_db/models.go
	rm -f internal/work_db/work_db/query.sql.go
	rm -f internal/search_db/search_db/db.go
	rm -f internal/search_db/search_db/models.go
	rm -f internal/search_db/search_db/query.sql.go
	rm -f cmd/*/service.exe

.PHONY: generate
generate:
	cd internal/sqlite ; make generate
	cd internal/work_db ; make generate
	cd internal/search_db ; make generate

.PHONY: config
config:
	cd config ; make all || exit 1

docker: 
	docker build -t jemison/dev -f Dockerfile.dev .
	docker build -t jemison/build -f Dockerfile.build .

.PHONY: build
# lint
build: clean config generate 
	cd cmd/migrate ; make build
	cd cmd/admin ; make build
	cd cmd/entree ; make build
	cd cmd/extract ; make build
	cd cmd/fetch ; make build
	cd cmd/pack ; make build
	cd cmd/serve ; make build
	cd cmd/validate ; make build
	cd cmd/walk ; make build
	cd assets ; rm -rf static/assets ; unzip -qq static.zip

.PHONY: up
up: build
	docker compose up

multi: build
	docker compose -f multi.yaml up \
		--scale entree=1 \
		--scale fetch=1 \
		--scale extract=4 \
		--scale walk=4


.PHONY: backend
backend: 
	docker compose -f backend.yaml up

.PHONY: 
macup:
	docker run -v ${PWD}:/app -t jemison/build
	docker compose up

.PHONY: run
run: clean generate
	cd assets ; unzip -qq -o static.zip > /dev/null 2>&1
	docker compose up

.PHONY: cloc
cloc:
	docker run --rm -v ${PWD}:/tmp aldanial/cloc --exclude-dir=assets .

delete_admin:
	cf delete -f admin

delete_entree:
	cf delete -f entree

delete_extract:
	cf delete -f extract

delete_fetch:
	cf delete -f fetch

delete_pack:
	cf delete -f pack

delete_serve:
	cf delete -f serve

delete_walk:
	cf delete -f walk

.PHONY: delete_all
delete_all: delete_admin delete_entree delete_extract delete_fetch delete_pack delete_serve delete_walk


# I need to delete_all every time, because there is not enough RAM
# in the sandbox to rolling deploy
.PHONY: terraform
terraform: delete_all
	docker run -v ${PWD}:/app -t jemison/build
	cd terraform ; make apply_all

.PHONY: docker_full_clean
docker_full_clean:
	-docker stop $(docker ps -a -q)
	-docker rm $(docker ps -a -q)

.PHONY: lint
lint:
	-golangci-lint run -v