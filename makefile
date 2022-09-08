GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

MAKEFLAGS += --silent

build:
	go build -o terraform-w-keycloak

build-example: build
	mkdir -p example/.terraform/plugins/terraform.local/charlesderek/keycloak/3.0.0/darwin_amd64
	mkdir -p example/terraform.d/plugins/terraform.local/charlesderek/keycloak/3.0.0/darwin_amd64
	cp terraform-w-keycloak example/.terraform/plugins/terraform.local/charlesderek/keycloak/3.0.0/darwin_amd64/
	cp terraform-w-keycloak example/terraform.d/plugins/terraform.local/charlesderek/keycloak/3.0.0/darwin_amd64/

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
	go test -v github.com/charlesderek/terraform-w-keycloak/keycloak
	TF_ACC=1 CHECKPOINT_DISABLE=1 go test -v -timeout 60m -count=1 github.com/charlesderek/terraform-w-keycloak/provider $(TESTARGS)

fmtcheck:
	lineCount=$(shell gofmt -l -s $(GOFMT_FILES) | wc -l | tr -d ' ') && exit $$lineCount

vet:
	go vet ./...

user-federation-example:
	cd custom-user-federation-example && ./gradlew shadowJar
