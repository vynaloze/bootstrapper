package main

import (
	"bootstrapper/actor/git"
	"bootstrapper/actor/terraform"
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
	tfcToken, ok := secrets.Get("actor.terraform.tfc.token")
	if !ok {
		panic(err)
	}

	commonGitOpts := git.Opts{
		Provider: "github.com",
		Project:  "bootstrapper-demo-org",

		RemoteAuthUser: "bootstrapper-demo",
		RemoteAuthPass: ghToken,
	}
	sharedInfraGitOpts, cicdRepoOpts := commonGitOpts, commonGitOpts
	sharedInfraGitOpts.Repo = "tf-infra-shared"
	cicdRepoOpts.Repo = "cicd"

	terraformOpts := blueprint.TerraformOpts{
		Opts: terraform.Opts{
			TerraformCloudOrg:   "bootstrapper-demo",
			TerraformCloudToken: tfcToken,
			TfVars: map[string]string{
				"repo_password": ghToken,
			},
		},
	}

	opts := blueprint.BootstrapOpts{
		SharedInfraRepoOpts: sharedInfraGitOpts,
		CICDRepoOpts:        cicdRepoOpts,
		TerraformOpts:       terraformOpts,
	}

	err = blueprint.Bootstrap(&opts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()

	opts2 := blueprint.SetupCICDRepoOpts{
		CICDRepoOpts: cicdRepoOpts,
	}

	err = blueprint.SetupCICDRepo(&opts2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press Enter to proceed")
	fmt.Scanln()

	opts3 := blueprint.AddCICDToRepoOpts{
		TargetRepoOpts: sharedInfraGitOpts,
		Templates: []blueprint.Template{
			{
				SourceFile: fmt.Sprintf("%s/terraform_infra.yml", sharedInfraGitOpts.Provider),
				Data: template.TerraformInfraTemplate{
					Project:         cicdRepoOpts.Project,
					Repo:            cicdRepoOpts.Repo,
					Module:          terraformOpts.GetTerraformInfraCoreDir(),
					Workflow:        "ci",
					DefaultBranch:   cicdRepoOpts.GetDefaultBranch(),
					OnDefaultBranch: false,
				},
				TargetFile: fmt.Sprintf(".github/workflows/%s.ci.yml", terraformOpts.GetTerraformInfraCoreDir()),
			},
			{
				SourceFile: fmt.Sprintf("%s/terraform_infra.yml", sharedInfraGitOpts.Provider),
				Data: template.TerraformInfraTemplate{
					Project:         cicdRepoOpts.Project,
					Repo:            cicdRepoOpts.Repo,
					Module:          terraformOpts.GetTerraformInfraCoreDir(),
					Workflow:        "cd",
					DefaultBranch:   cicdRepoOpts.GetDefaultBranch(),
					OnDefaultBranch: true,
				},
				TargetFile: fmt.Sprintf(".github/workflows/%s.cd.yml", terraformOpts.GetTerraformInfraCoreDir()),
			},
			{
				SourceFile: ".tflint.hcl",
				Data:       nil,
				TargetFile: ".tflint.hcl",
			},
		},
	}

	err = blueprint.AddCICDToRepo(&opts3)
	if err != nil {
		log.Fatalln(err)
	}
}
