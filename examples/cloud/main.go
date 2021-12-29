package main

import (
	"bootstrapper/actor/git"
	"bootstrapper/blueprint"
	"bootstrapper/datasource"
	"bootstrapper/template"
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

	//sharedInfraGitOpts := git.Opts{
	//	Provider: gitProvider, Project: gitProject, Repo: "tf-infra-shared",
	//	RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	//}
	cicdRepoOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "cicd",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}
	tfEnvRepoOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "tf-env",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}

	log.Printf("start cloud phase")
	log.Printf("setup environment module")
	//createGitRepoOpts := blueprint.CreateGitRepoOpts{
	//	SharedInfraRepoOpts: sharedInfraGitOpts,
	//	NewRepoOpts:         tfEnvRepoOpts,
	//	NewRepoType:         template.TerraformModule,
	//}
	//err = blueprint.CreateGitRepo(createGitRepoOpts)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//fmt.Println("Press Enter to proceed")
	//fmt.Scanln()

	setupCloudEnvOpts := blueprint.SetupCloudEnvModuleOpts{
		EnvRepoOpts:   tfEnvRepoOpts,
		CloudProvider: cloudProvider,
		CICDTemplates: blueprint.TfEnvCICDPreset(tfEnvRepoOpts, cicdRepoOpts),
	}

	err = blueprint.SetupCloudEnvModule(&setupCloudEnvOpts)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println("Press Enter to proceed")
	//fmt.Scanln()

	log.Printf("setup infra module(s)")

	// TODO cloud infra env here

}
