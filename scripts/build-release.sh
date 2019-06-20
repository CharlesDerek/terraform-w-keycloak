#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

mkdir -p ../artifacts

for config in $(cat release-targets.json | jq -rc '.[]'); do
	os=$(echo ${config} | jq -r '.os')
	platform=$(echo ${config} | jq -r '.platform')
	static=$(echo ${config} | jq -r '.static // false')
  linkage=''
  if [[ ${static} = 'true' ]]; then
    export CGO_ENABLED=0
    linkage='_static'
  else
    unset CGO_ENABLED
  fi

	echo "Building for ${os}_${platform}${linkage}..."

	GOOS=${os} GOARCH=${platform} go build -o terraform-w-keycloak_v${CIRCLE_TAG} ..
	zip terraform-w-keycloak_v${CIRCLE_TAG}_${os}_${platform}${linkage}.zip terraform-w-keycloak_v${CIRCLE_TAG} ../LICENSE
	mv terraform-w-keycloak_v${CIRCLE_TAG}_${os}_${platform}${linkage}.zip ../artifacts
	rm terraform-w-keycloak_v${CIRCLE_TAG}
done;

cd ../artifacts

sha256sum -b * > SHA256SUMS
