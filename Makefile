.PHONY: clean
clean:
	rm -f internal/postgres/work_db/schema/db.go
	rm -f internal/postgres/work_db/schema/models.go
	rm -f internal/postgres/work_db/schema/query.sql.go
	rm -f internal/postgres/search_db/schema/db.go
	rm -f internal/postgres/search_db/schema/models.go
	rm -f internal/postgres/search_db/schema/query.sql.go
	rm -f cmd/*/service.exe

.PHONY: generate
generate:
	cd internal/postgres ; make generate
	# cd internal/postgres/search_db ; make generate

.PHONY: config
config:
	cd config ; make all || exit 1

docker: 
	docker build -t jemison/dev -f Dockerfile.dev .
	docker build -t jemison/build -f Dockerfile.build .

.PHONY: build
# lint
build: clean config generate 
	echo "build migrate"
	cd cmd/migrate ; make build
	echo "build admin"
	cd cmd/admin ; make build
	echo "build entree"
	cd cmd/entree ; make build
	echo "build extract"
	cd cmd/extract ; make build
	echo "build fetch"
	cd cmd/fetch ; make build
	echo "build collect"
	cd cmd/collect ; make build
	echo "build pack"
	cd cmd/pack ; make build
	echo "build serve"
	cd cmd/serve ; make build
	echo "build validate"
	cd cmd/validate ; make build
	echo "build walk"
	cd cmd/walk ; make build

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
	docker compose up

.PHONY: cloc
cloc:
	docker run --rm -v ${PWD}:/tmp aldanial/cloc --exclude-dir=assets .

.PHONY: delete_all
delete_all: 
	cd terraform ; make cfclean

# I need to delete_all every time, because there is not enough RAM
# in the sandbox to rolling deploy
.PHONY: terraform
terraform: delete_all build
	docker run -v ${PWD}:/app -t jemison/build
	cd terraform ; make apply_all

.PHONY: docker_full_clean
docker_full_clean:
	-docker stop $(docker ps -a -q)
	-docker rm $(docker ps -a -q)

.PHONY: lint
lint:
	-golangci-lint run -v