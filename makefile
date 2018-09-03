TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-w-keycloak

example: build
	mkdir -p example/.terraform/plugins/darwin_amd64
	cp terraform-w-keycloak example/.terraform/plugins/darwin_amd64/

local: deps
	docker-compose up --build -d
	./scripts/wait-for-local-keycloak.sh
	./scripts/create-terraform-client.sh

deps:
	./scripts/check-deps.sh

fmt:
	gofmt -w -s $(GOFMT_FILES)

test: fmtcheck vet
	go test $(TEST)

testacc: fmtcheck vet
	TF_ACC=1 go test $(TEST) -v $(TESTARGS)

fmtcheck:
	lineCount=$(shell gofmt -l -s $(GOFMT_FILES) | wc -l | tr -d ' ') && exit $$lineCount

vet:
	go vet ./...
