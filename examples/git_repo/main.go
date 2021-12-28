package main

import (
	"bootstrapper/actor/git"
	"bootstrapper/blueprint"
	"bootstrapper/datasource"
	"bootstrapper/template"
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

	commonGitOpts := git.Opts{
		Provider: "github.com",
		Project:  "bootstrapper-demo-org",

		RemoteAuthUser: "bootstrapper-demo",
		RemoteAuthPass: ghToken,
	}
	sharedInfraGitOpts, newRepoOpts := commonGitOpts, commonGitOpts
	sharedInfraGitOpts.Repo = "tf-infra-shared"
	newRepoOpts.Repo = "tf-env"

	opts := blueprint.CreateGitRepoOpts{
		SharedInfraRepoOpts: sharedInfraGitOpts,
		NewRepoOpts:         newRepoOpts,
		NewRepoType:         template.TerraformModule,
	}

	err = blueprint.CreateGitRepo(opts)
	if err != nil {
		fmt.Println(err)
	}
}
