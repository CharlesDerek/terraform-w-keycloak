GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
GOOS?=darwin
GOARCH?=amd64

MAKEFLAGS += --silent

VERSION=$$(git describe --tags)

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o terraform-w-keycloak_$(VERSION)

build-example: build
	mkdir -p example/.terraform/plugins/terraform.local/charlesderek/keycloak/4.0.0/$(GOOS)_$(GOARCH)
	mkdir -p example/terraform.d/plugins/terraform.local/charlesderek/keycloak/4.0.0/$(GOOS)_$(GOARCH)
	cp terraform-w-keycloak example/.terraform/plugins/terraform.local/charlesderek/keycloak/4.0.0/$(GOOS)_$(GOARCH)/
	cp terraform-w-keycloak example/terraform.d/plugins/terraform.local/charlesderek/keycloak/4.0.0/$(GOOS)_$(GOARCH)/

local: deps
	docker compose up --build -d
	./scripts/wait-for-local-keycloak.sh
	./scripts/create-terraform-client.sh

deps:
	./scripts/check-deps.sh

fmt:
	gofmt -w -s $(GOFMT_FILES)

test: fmtcheck vet
	go test $(TEST)

testacc: fmtcheck vet
	go test -v github.com/charlesderek/terraform-w-keycloak/keycloak
	TF_ACC=1 CHECKPOINT_DISABLE=1 go test -v -timeout 60m -parallel 4 github.com/charlesderek/terraform-w-keycloak/provider $(TESTARGS)

fmtcheck:
	lineCount=$(shell gofmt -l -s $(GOFMT_FILES) | wc -l | tr -d ' ') && exit $$lineCount

vet:
	go vet ./...

user-federation-example:
	cd custom-user-federation-example && ./gradlew shadowJar
