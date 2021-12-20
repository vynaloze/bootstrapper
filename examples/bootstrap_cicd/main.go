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

	cicdRepoOpts := git.Opts{
		Provider: "github.com",
		Project:  "bootstrapper-demo-org",
		Repo:     "cicd",

		RemoteAuthUser: "bootstrapper-demo",
		RemoteAuthPass: ghToken,
	}

	opts := blueprint.SetupCICDRepoOpts{
		CICDRepoOpts: cicdRepoOpts,
	}

	err = blueprint.SetupCICDRepo(&opts)
	if err != nil {
		fmt.Println(err)
	}
}
