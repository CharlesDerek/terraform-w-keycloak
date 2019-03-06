#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

mkdir -p ../artifacts

for config in $(cat release-targets.json | jq -rc '.[]'); do
	os=$(echo ${config} | jq -r '.os')
	platform=$(echo ${config} | jq -r '.platform')

	echo "Building for ${os}_${platform}..."

	GOOS=${os} GOARCH=${platform} go build -o terraform-w-keycloak_v${CIRCLE_TAG} ..
	zip terraform-w-keycloak_v${CIRCLE_TAG}_${os}_${platform}.zip terraform-w-keycloak_v${CIRCLE_TAG} ../LICENSE
	mv terraform-w-keycloak_v${CIRCLE_TAG}_${os}_${platform}.zip ../artifacts
	rm terraform-w-keycloak_v${CIRCLE_TAG}
done;

cd ../artifacts

sha256sum -b * > SHA256SUMS
