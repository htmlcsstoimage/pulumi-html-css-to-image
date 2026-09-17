VERSION ?= $(shell cat VERSION)

.PHONY: build schema test vet fmt sdk
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
	pulumi package gen-sdk internal/provider/schema.json --language nodejs --out sdk
	pulumi package gen-sdk internal/provider/schema.json --language python --out sdk
	pulumi package gen-sdk internal/provider/schema.json --language go --out sdk
	test -f sdk/go/go.mod || go -C sdk/go mod init github.com/htmlcsstoimage/pulumi-html-css-to-image/sdk/go
	go -C sdk/go mod edit -go=1.25.11 -require=github.com/pulumi/pulumi/sdk/v3@v3.259.0
	cd sdk/go && GOWORK=off go mod tidy
	pulumi package gen-sdk internal/provider/schema.json --language dotnet --out sdk
	pulumi package gen-sdk internal/provider/schema.json --language java --out sdk
	node scripts/prepare-java-sdk.mjs
