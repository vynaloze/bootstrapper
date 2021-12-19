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
	ghToken, ok := secrets.Get("actor.git.github.token")
	if !ok {
		panic(err)
	}
	//tfcToken, ok := secrets.Get("actor.terraform.tfc.token")
	//if !ok {
	//	panic(err)
	//}

	commonGitOpts := git.Opts{
		Provider: "github.com",
		Project:  "bootstrapper-demo-org",

		RemoteAuthUser: "bootstrapper-demo",
		RemoteAuthPass: ghToken,
	}
	sharedInfraGitOpts, cicdRepoOpts := commonGitOpts, commonGitOpts
	sharedInfraGitOpts.Repo = "tf-infra-shared"
	cicdRepoOpts.Repo = "cicd"

	opts := blueprint.SetupCICDRepoOpts{
		CICDRepoOpts: cicdRepoOpts,
	}

	err = blueprint.SetupCICDRepo(&opts)
	if err != nil {
		fmt.Println(err)
	}
}
