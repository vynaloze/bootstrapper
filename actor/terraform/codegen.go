package terraform

import (
	hclencoder_blocks "github.com/rodaine/hclencoder"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

type providersTf struct {
	Github *providersTfGithub `hcl:"provider,block"`
	Tfe    providersTfTfe     `hcl:"provider,block"`
}

type providersTfGithub struct {
	Provider string `hcl:",key"`
	Owner    string `hcl:"owner"`
	Token    string `hcl:"token"`
}

type providersTfTfe struct {
	Provider string `hcl:",key"`
	Token    string `hcl:"token"`
}

func (a *Actor) writeProvidersTf(dir string) error {
	providersTf := providersTf{
		Tfe: providersTfTfe{
			Provider: "tfe",
			Token:    a.opts.ProviderSecrets["tfe"]["token"],
		},
	}
	if _, ok := a.opts.ProviderSecrets["github"]; ok {
		providersTf.Github = &providersTfGithub{
			Provider: "github",
			Owner:    a.opts.ProviderSecrets["github"]["owner"],
			Token:    a.opts.ProviderSecrets["github"]["token"],
		}
	}
	providersTfContent, err := hclencoder_blocks.Encode(providersTf)
	if err != nil {
		return err
	}
	providersTfFile := filepath.Join(dir, "providers_override.tf")
	return ioutil.WriteFile(providersTfFile, providersTfContent, 0644)
}
