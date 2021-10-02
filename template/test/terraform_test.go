package test

import (
	"bootstrapper/template"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerraformVersionsFile(t *testing.T) {
	providers := map[string]template.TerraformProvider{
		"ourcloud": {
			Source:  "terraform.example.com/examplecorp/ourcloud",
			Version: "~> 1.0",
		},
		"aws": {
			Source:  "hashicorp/aws",
			Version: ">= 2.7.0",
		},
		"mycorp-http": {
			Source:  "mycorp/http",
			Version: "!= 1.0, < 2.0.0",
		},
		"random_2317": {
			Source:  "kw/random2317",
			Version: "= 1.2.0-beta",
		},
	}
	expectedFile := `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 2.7.0"
    }
    mycorp-http = {
      source  = "mycorp/http"
      version = "!= 1.0, < 2.0.0"
    }
    ourcloud = {
      source  = "terraform.example.com/examplecorp/ourcloud"
      version = "~> 1.0"
    }
    random_2317 = {
      source  = "kw/random2317"
      version = "= 1.2.0-beta"
    }
  }
}
`
	actualFile, err := template.TerraformVersionsFile(providers)
	assert.Nil(t, err)
	assert.Equal(t, expectedFile, actualFile)
}
