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

	commonGitOpts := git.Opts{
		Provider: "github.com",
		Project:  "bootstrapper-demo-org",

		RemoteAuthUser: "bootstrapper-demo",
		RemoteAuthPass: ghToken,
	}
	sharedInfraGitOpts, newRepoOpts, cicdRepoOpts := commonGitOpts, commonGitOpts, commonGitOpts
	sharedInfraGitOpts.Repo = "tf-infra-shared"
	newRepoOpts.Repo = "tf-env"
	cicdRepoOpts.Repo = "cicd"

	opts := blueprint.CreateGitRepoOpts{
		SharedInfraRepoOpts: sharedInfraGitOpts,
		NewRepoOpts:         newRepoOpts,
		NewRepoType:         template.TerraformModule,
	}

	err = blueprint.CreateGitRepo(opts)
	if err != nil {
		fmt.Println(err)
	}

	cicdOpts := blueprint.AddCICDToRepoOpts{ //TODO don't worry about it - will replace with app repo sometime at the end
		TargetRepoOpts: newRepoOpts,
		Templates: []blueprint.Template{
			{
				SourceFile: fmt.Sprintf("%s/%s_ci.yml", newRepoOpts.Provider, template.TerraformModule),
				Data: template.TerraformModuleTemplate{
					Project:       cicdRepoOpts.Project,
					Repo:          cicdRepoOpts.Repo,
					Modules:       []string{"base", "k8s"},
					DefaultBranch: cicdRepoOpts.GetDefaultBranch(),
				},
				TargetFile: ".github/workflows/ci.yml",
			},
			{
				SourceFile: fmt.Sprintf("%s/%s_cd.yml", newRepoOpts.Provider, template.TerraformModule),
				Data: template.TerraformModuleTemplate{
					Project:       cicdRepoOpts.Project,
					Repo:          cicdRepoOpts.Repo,
					Modules:       []string{"base", "k8s"},
					DefaultBranch: cicdRepoOpts.GetDefaultBranch(),
				},
				TargetFile: ".github/workflows/cd.yml",
			},
			{
				SourceFile: ".tflint.hcl",
				Data:       nil,
				TargetFile: ".tflint.hcl",
			},
			{
				SourceFile: ".releaserc.yaml",
				Data:       nil,
				TargetFile: ".releaserc.yaml",
			},
		},
	}

	err = blueprint.AddCICDToRepo(&cicdOpts)
	if err != nil {
		log.Fatalln(err)
	}
}
