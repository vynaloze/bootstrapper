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

	opts := blueprint.ApplicationGitRepoOpts{
		RemoteOpts: git.RemoteOpts{
			URL:  "github.com/bootstrapper-demo",
			Auth: token,
		},
		RepoName: "after_refactor",
	}

	err = blueprint.CreateApplicationGitRepo(opts)
	if err != nil {
		fmt.Println(err)
	}
}
