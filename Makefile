VERSION ?= $(shell cat VERSION)

.PHONY: build schema test vet fmt sdk sdk-nodejs sdk-python sdk-go sdk-dotnet sdk-java
schema:
	go run ./cmd/schema -version $(VERSION)
build: schema
	go build -ldflags "-X main.version=$(VERSION)" -o bin/pulumi-resource-html-css-to-image ./cmd/pulumi-resource-html-css-to-image
test:
	go test -race ./...
vet:
	go vet ./...
fmt:
	gofmt -w cmd internal
sdk: schema
	$(MAKE) sdk-nodejs sdk-python sdk-go sdk-dotnet sdk-java

sdk-nodejs:
	pulumi package gen-sdk internal/provider/schema.json --language nodejs --out sdk
	node scripts/prepare-sdk-readme.mjs nodejs
sdk-python:
	pulumi package gen-sdk internal/provider/schema.json --language python --out sdk
	node scripts/prepare-sdk-readme.mjs python
sdk-go:
	pulumi package gen-sdk internal/provider/schema.json --language go --out sdk
	node scripts/prepare-sdk-readme.mjs go
	test -f sdk/go/go.mod || go -C sdk/go mod init github.com/htmlcsstoimage/pulumi-html-css-to-image/sdk/go
	go -C sdk/go mod edit -go=1.25.11 -require=github.com/pulumi/pulumi/sdk/v3@v3.259.0
	cd sdk/go && GOWORK=off go mod tidy
sdk-dotnet:
	pulumi package gen-sdk internal/provider/schema.json --language dotnet --out sdk
	node scripts/prepare-sdk-readme.mjs dotnet
sdk-java:
	pulumi package gen-sdk internal/provider/schema.json --language java --out sdk
	node scripts/prepare-java-sdk.mjs
	node scripts/prepare-sdk-readme.mjs java
