package main

import (
	"bootstrapper/actor/git"
	"bootstrapper/blueprint"
	"bootstrapper/datasource"
	"bootstrapper/template"
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

	gitProvider := "github.com"
	gitProject := "bootstrapper-demo-org"
	gitUser := "bootstrapper-demo"
	gitPass := ghToken

	cloudProvider := template.AWS

	env := "stg"

	sharedInfraGitOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "tf-infra-shared",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}
	cicdRepoOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "cicd",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}
	tfEnvRepoOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "tf-env",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}
	tfInfraRepoOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "tf-infra-" + env,
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}

	log.Printf("start cloud phase")

	log.Printf("setup environment module")
	createTfEnvRepo(sharedInfraGitOpts, tfEnvRepoOpts)
	setupTfEnvRepo(tfEnvRepoOpts, cloudProvider, cicdRepoOpts)

	log.Printf("setup infra module(s)")
	createTfInfraRepo(sharedInfraGitOpts, tfInfraRepoOpts)

}

func createTfEnvRepo(sharedInfraGitOpts git.Opts, tfEnvRepoOpts git.Opts) {
	createGitRepoOpts := blueprint.CreateGitRepoOpts{
		SharedInfraRepoOpts: sharedInfraGitOpts,
		NewRepoOpts:         tfEnvRepoOpts,
		NewRepoType:         template.TerraformModule,
	}
	err := blueprint.CreateGitRepo(createGitRepoOpts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()
}

func setupTfEnvRepo(tfEnvRepoOpts git.Opts, cloudProvider template.CloudProvider, cicdRepoOpts git.Opts) {
	setupCloudEnvOpts := blueprint.SetupCloudEnvModuleOpts{
		EnvRepoOpts:   tfEnvRepoOpts,
		CloudProvider: cloudProvider,
		CICDTemplates: blueprint.TfEnvCICDPreset(tfEnvRepoOpts, cicdRepoOpts),
	}

	err := blueprint.SetupCloudEnvModule(&setupCloudEnvOpts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()
}

func createTfInfraRepo(sharedInfraGitOpts git.Opts, tfInfraRepoOpts git.Opts) {
	createGitRepoOpts := blueprint.CreateGitRepoOpts{
		SharedInfraRepoOpts: sharedInfraGitOpts,
		NewRepoOpts:         tfInfraRepoOpts,
		NewRepoType:         template.TerraformInfra,
		NewRepoExtraContent: template.GitRepoExtraContent{Modules: []string{"base", "k8s"}},
	}
	err := blueprint.CreateGitRepo(createGitRepoOpts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()
}
