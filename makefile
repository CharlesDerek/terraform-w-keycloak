all=build

MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-w-keycloak
