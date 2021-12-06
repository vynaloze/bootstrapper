package main

import (
	"bootstrapper/actor/git"
	"bootstrapper/blueprint"
	"bootstrapper/datasource"
	"fmt"
)

func main() {
	secrets, err := datasource.NewYamlFile("secrets.yaml")
	if err != nil {
		panic(err)
	}
	token, ok := secrets.Get("actor.git.github.token")
	if !ok {
		panic(err)
	}

	opts := blueprint.BootstrapOpts{
		Opts: git.Opts{
			RemoteBaseURL:  "https://github.com/bootstrapper-demo",
			RemoteAuthUser: "bootstrapper-demo",
			RemoteAuthPass: token,
		},
		TerraformOpts: blueprint.TerraformOpts{
			ProviderSecrets: map[string]string{"GITHUB_TOKEN": token},
		},
	}

	err = blueprint.Bootstrap(&opts)
	if err != nil {
		fmt.Println(err)
	}
}
