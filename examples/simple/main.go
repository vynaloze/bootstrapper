package main

import (
	"bootstrapper/blueprint"
	"bootstrapper/datasource"
	"fmt"
)

func main() {
	l, _ := datasource.NewLiteral()
	l["git.provider"] = "github.com"
	l["git.project"] = "bootstrapper-demo"

	_, err := datasource.NewYamlFile("secrets.yaml")
	if err != nil {
		fmt.Println(err)
	}

	err = blueprint.CreateApplicationGitRepo("third_app")
	if err != nil {
		fmt.Println(err)
	}
}
