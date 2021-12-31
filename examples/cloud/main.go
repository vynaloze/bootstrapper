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
	terraformCloudOrg := "bootstrapper-demo"

	region := "eu-central-1"
	env := "stg"

	baseVars := template.TfInfraBaseTfVars{
		Region:              region,
		Environment:         env,
		BaseDomain:          "b-demo.org",
		VpcCidr:             "172.20.0.0/16",
		ClientCidrBlock:     "172.30.0.0/22",
		PrivateSubnetsCidrs: []string{"172.20.0.0/20", "172.20.16.0/20"},
		PublicSubnetsCidrs:  []string{"172.20.128.0/20", "172.20.144.0/20"},
	}
	k8sVars := template.TfInfraK8sTfVars{
		Region:              region,
		Environment:         env,
		VpcCidr:             "172.21.0.0/16",
		PrivateSubnetsCidrs: []string{"172.21.0.0/20", "172.21.16.0/20"},
		PublicSubnetsCidrs:  []string{"172.21.128.0/20", "172.21.144.0/20"},
	}

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
	setupTfInfraRepoBase(tfInfraRepoOpts, tfEnvRepoOpts, cicdRepoOpts, cloudProvider, terraformCloudOrg, baseVars)
	setupTfInfraRepoK8s(tfInfraRepoOpts, tfEnvRepoOpts, cicdRepoOpts, cloudProvider, terraformCloudOrg, k8sVars)

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

func setupTfInfraRepoBase(tfInfraRepoOpts git.Opts, tfEnvRepoOpts git.Opts, cicdRepoOpts git.Opts,
	cloudProvider template.CloudProvider, terraformCloudOrg string, baseVars template.TfInfraBaseTfVars) {
	setupCloudInfraOpts := blueprint.SetupCloudInfraOpts{
		InfraRepoOpts: tfInfraRepoOpts,
		EnvRepoOpts:   tfEnvRepoOpts,

		CloudProvider:     cloudProvider,
		TerraformCloudOrg: terraformCloudOrg,

		CICDTemplates: blueprint.TfInfraCICDPreset(tfInfraRepoOpts, cicdRepoOpts,
			[]template.CICDTerraformInfraModuleTemplate{
				{Name: "base"},
			}),
		TerraformVars: baseVars,

		Module: "base",
	}

	err := blueprint.SetupCloudInfra(&setupCloudInfraOpts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()
}

func setupTfInfraRepoK8s(tfInfraRepoOpts git.Opts, tfEnvRepoOpts git.Opts, cicdRepoOpts git.Opts,
	cloudProvider template.CloudProvider, terraformCloudOrg string, k8sVars template.TfInfraK8sTfVars) {
	setupCloudInfraOpts := blueprint.SetupCloudInfraOpts{
		InfraRepoOpts: tfInfraRepoOpts,
		EnvRepoOpts:   tfEnvRepoOpts,

		CloudProvider:     cloudProvider,
		TerraformCloudOrg: terraformCloudOrg,

		CICDTemplates: blueprint.TfInfraCICDPreset(tfInfraRepoOpts, cicdRepoOpts,
			[]template.CICDTerraformInfraModuleTemplate{
				{Name: "base"},
				{Name: "k8s", Dependencies: []string{"base"}},
			}),
		TerraformVars:          k8sVars,
		TerraformParentModules: []string{"base"},

		Module: "k8s",
	}

	err := blueprint.SetupCloudInfra(&setupCloudInfraOpts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()
}
