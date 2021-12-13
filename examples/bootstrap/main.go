package main

import (
	"bootstrapper/actor/git"
	"bootstrapper/actor/terraform"
	"bootstrapper/blueprint"
	"bootstrapper/datasource"
	"fmt"
)

func main() {
	secrets, err := datasource.NewYamlFile("secrets.yaml")
	if err != nil {
		panic(err)
	}
	ghToken, ok := secrets.Get("actor.git.github.token")
	if !ok {
		panic(err)
	}
	tfcToken, ok := secrets.Get("actor.terraform.tfc.token")
	if !ok {
		panic(err)
	}

	opts := blueprint.BootstrapOpts{
		SharedInfraRepoOpts: git.Opts{
			Provider: "github.com",
			Project:  "bootstrapper-demo-org",
			Repo:     "tf-infra-shared",

			RemoteAuthUser: "bootstrapper-demo",
			RemoteAuthPass: ghToken,
		},
		TerraformOpts: blueprint.TerraformOpts{
			Opts: terraform.Opts{
				TerraformCloudOrg:   "bootstrapper-demo",
				TerraformCloudToken: tfcToken,
				TfVars: map[string]string{
					"repo_password": ghToken,
				},
			},
		},
	}

	err = blueprint.Bootstrap(&opts)
	if err != nil {
		fmt.Println(err)
	}
}
