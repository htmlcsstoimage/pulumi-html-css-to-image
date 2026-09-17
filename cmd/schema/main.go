// Command schema generates the Pulumi package from Terraform's schemas.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/htmlcsstoimage/pulumi-html-css-to-image/internal/provider"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfgen"
)

func main() {
	version := flag.String("version", "0.1.0", "Provider version")
	flag.Parse()
	// The bridge's default docs lookup assumes a checkout directory named
	// "provider". Resolve our pinned module explicitly so local and CI generation
	// read the same documentation regardless of the checkout's name.
	lookup := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/htmlcsstoimage/terraform-provider-html-css-to-image")
	lookup.Env = append(os.Environ(), "GOWORK=off")
	moduleDir, err := lookup.Output()
	if err != nil {
		log.Fatalf("resolve Terraform provider documentation: %v", err)
	}
	info := provider.Info(*version)
	info.UpstreamRepoPath = strings.TrimSpace(string(moduleDir))
	if info.UpstreamRepoPath == "" {
		log.Fatal("Terraform provider module directory is empty")
	}
	generated, err := tfgen.GenerateSchema(context.Background(), tfgen.GenerateSchemaOptions{ProviderInfo: info, XInMemoryDocs: true})
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
