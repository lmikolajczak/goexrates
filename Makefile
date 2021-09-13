## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}


## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database postgres://admin:admin@db/exrates?sslmode=disable up


## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Running down migrations...'
	migrate -path ./migrations -database postgres://admin:admin@db/exrates?sslmode=disable down


## api/build
.PHONY: api/build
api/build:
	@echo 'Building api...'
	@cd cmd/api && go build -o /go/bin/goexrates-api


## cli/build
.PHONY: cli/build
cli/build:
	@echo 'Building cli...'
	@cd cmd/cli && go build -o /go/bin/goexrates-cli
