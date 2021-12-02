package main

import (
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
		EnvVars: map[string]string{"GITHUB_TOKEN": token},
	}

	err = blueprint.Bootstrap(opts)
	if err != nil {
		fmt.Println(err)
	}
}
