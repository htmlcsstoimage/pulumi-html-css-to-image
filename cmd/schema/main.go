// Command schema generates the Pulumi package from Terraform's schemas.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/htmlcsstoimage/pulumi-html-css-to-image/internal/provider"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfgen"
)

func main() {
	version := flag.String("version", "0.1.0", "Provider version")
	flag.Parse()
	generated, err := tfgen.GenerateSchema(context.Background(), tfgen.GenerateSchemaOptions{ProviderInfo: provider.Info(*version), XInMemoryDocs: true})
	if err != nil {
		log.Fatal(err)
	}
	var value any
	if err = json.Unmarshal(generated.ProviderMetadata.PackageSchema, &value); err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile("internal/provider/schema.json", append(data, '\n'), 0644); err != nil {
		log.Fatal(err)
	}
}
