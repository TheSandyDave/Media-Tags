package main

import _ "embed"

//go:embed generated/api/openapi.json
var spec []byte

//go:generate curl  --create-dirs --output ./tmp/openapi-generator-cli.jar -sSLO https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.9.0/openapi-generator-cli-7.9.0.jar
//go:generate -command apigen java -Dlog.level=off -jar ./tmp/openapi-generator-cli.jar generate -i docs/openapi.yaml
//go:generate apigen -g go-gin-server -t ./template -s --minimal-update --remove-operation-id-prefix --skip-operation-example --additional-properties=packageName=restgen,apiPath=generated/api
//go:generate apigen -g openapi -o ./generated/api --additional-properties=outputFileName=openapi.json
//go:generate go fmt ./generated/api
