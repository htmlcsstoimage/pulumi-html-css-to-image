package main

import (
	"context"

	"github.com/htmlcsstoimage/pulumi-html-css-to-image/internal/provider"
	pf "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
)

var version = "0.1.2"

func main() {
	pf.Main(context.Background(), provider.Name, provider.Info(version), pf.ProviderMetadata{PackageSchema: provider.Schema})
}
