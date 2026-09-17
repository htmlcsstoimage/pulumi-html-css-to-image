package provider

import _ "embed"

// Schema is generated from the Terraform schemas by make schema.
//
//go:embed schema.json
var Schema []byte
