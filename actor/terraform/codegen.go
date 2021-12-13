package terraform

import (
	"github.com/hashicorp/hcl2/hclwrite"
	hclencoder_blocks "github.com/rodaine/hclencoder"
	"github.com/zclconf/go-cty/cty"
	"io/ioutil"
	"os"
	"path/filepath"
)

const secretsFile = "secrets_override.auto.tfvars"

type terraformTfrc struct {
	Credentials terraformTfrcCredentials `hcl:"credentials,block"`
}

type terraformTfrcCredentials struct {
	Host  string `hcl:",key"`
	Token string `hcl:"token"`
}

func (a *Actor) writeCliConfigFile(dir string) error {
	terraformTfrc := terraformTfrc{
		Credentials: terraformTfrcCredentials{
			Host:  "app.terraform.io",
			Token: a.opts.TerraformCloudToken,
		},
	}
	terraformTfrcContent, err := hclencoder_blocks.Encode(terraformTfrc)
	if err != nil {
		return err
	}
	cliConfigFile := filepath.Join(dir, "terraform.tfrc")
	err = os.Setenv("TF_CLI_CONFIG_FILE", cliConfigFile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cliConfigFile, terraformTfrcContent, 0644)
}

func (a *Actor) writeTfVarsFile(dir string) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	for k, v := range a.opts.TfVars {
		rootBody.SetAttributeValue(k, cty.StringVal(v))
	}

	cliConfigFile := filepath.Join(dir, secretsFile)
	return ioutil.WriteFile(cliConfigFile, f.Bytes(), 0644)
}
