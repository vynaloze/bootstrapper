package main

import (
	"bootstrapper/actor/git"
	"bootstrapper/actor/terraform"
	"bootstrapper/blueprint"
	"bootstrapper/datasource"
	"fmt"
	"log"
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

	gitProvider := "github.com"
	gitProject := "bootstrapper-demo-org"
	gitUser := "bootstrapper-demo"
	gitPass := ghToken

	sharedInfraGitOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "tf-infra-shared",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}
	cicdRepoOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "cicd",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}
	terraformOpts := blueprint.TerraformOpts{
		Opts: terraform.Opts{
			TerraformCloudOrg:       "bootstrapper-demo",
			TerraformCloudToken:     tfcToken,
			TerraformCloudWorkspace: "tf-infra-shared-core",
			ProviderSecrets: map[string]map[string]string{
				"github": {"owner": gitProject, "token": gitPass},
				"tfe":    {"token": tfcToken},
			},
		},
	}

	bootstrapOpts := blueprint.BootstrapOpts{
		SharedInfraRepoOpts: sharedInfraGitOpts,
		CICDRepoOpts:        cicdRepoOpts,
		TerraformOpts:       terraformOpts,
	}

	err = blueprint.Bootstrap(&bootstrapOpts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()

	setupCICDOpts := blueprint.SetupCICDRepoOpts{
		CICDRepoOpts: cicdRepoOpts,
	}

	err = blueprint.SetupCICDRepo(&setupCICDOpts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()

	addCICDOpts := blueprint.TfInfraSharedCICDPreset(sharedInfraGitOpts, cicdRepoOpts, terraformOpts)

	err = blueprint.AddCICDToRepo(&addCICDOpts)
	if err != nil {
		log.Fatalln(err)
	}
}
