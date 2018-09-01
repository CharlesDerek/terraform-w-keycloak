all=build

MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-w-keycloak

example: build
	mkdir -p example/.terraform/plugins/darwin_amd64
	cp terraform-w-keycloak example/.terraform/plugins/darwin_amd64/
	cd example && terraform init && terraform plan
