package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"log"
	"time"
)

type AddCICDToRepoOpts struct {
	TargetRepoOpts git.Opts
	Templates      []Template
}

func AddCICDToRepo(opts *AddCICDToRepoOpts) error {
	log.Printf("adding CICD to %s repo", opts.TargetRepoOpts.Repo)

	localActor := git.NewLocal(&opts.TargetRepoOpts)
	remoteActor, err := git.NewRemote(&opts.TargetRepoOpts)
	if err != nil {
		return fmt.Errorf("cannot initialize remote Git actor: %w", err)
	}

	log.Printf("preparing CICD pipelines templates")

	gitFiles, err := templatesToGitFiles("cicd/pipeline_templates", opts.Templates)
	if err != nil {
		return fmt.Errorf("error preparing CICD templates: %w", err)
	}

	log.Printf("pushing changes to remote repository")

	message := "chore: add CI/CD pipelines templates"
	branch := fmt.Sprintf("%s/%d", opts.TargetRepoOpts.GetAuthorName(), time.Now().UnixMilli())
	err = commitAndPush(localActor, branch, message, gitFiles)
	if err != nil {
		return err
	}
	err = remoteActor.RequestReview(&branch, &message)
	if err != nil {
		return fmt.Errorf("error creating PR: %w", err)
	}

	return nil
}

func commitAndPush(localActor git.LocalActor, branch string, message string, gitFiles []git.File) error {
	err := localActor.CommitMany(branch, message, gitFiles...)
	if err != nil {
		return fmt.Errorf("error committing files: %w", err)
	}
	err = localActor.Push()
	if err != nil {
		return fmt.Errorf("error pushing changes: %w", err)
	}
	return nil
}

func TfInfraCICDPreset(infraGitOpts git.Opts, cicdRepoOpts git.Opts, modules []template.CICDTerraformInfraModuleTemplate) []Template {
	return []Template{
		{
			SourceFile: fmt.Sprintf("%s/%s_ci.yml", infraGitOpts.Provider, template.TerraformInfra),
			Data: template.CICDTerraformInfraTemplate{
				Project:       cicdRepoOpts.Project,
				Repo:          cicdRepoOpts.Repo,
				Modules:       modules,
				DefaultBranch: cicdRepoOpts.GetDefaultBranch(),
			},
			TargetFile: ".github/workflows/ci.yml",
		},
		{
			SourceFile: fmt.Sprintf("%s/%s_cd.yml", infraGitOpts.Provider, template.TerraformInfra),
			Data: template.CICDTerraformInfraTemplate{
				Project:       cicdRepoOpts.Project,
				Repo:          cicdRepoOpts.Repo,
				Modules:       modules,
				DefaultBranch: cicdRepoOpts.GetDefaultBranch(),
			},
			TargetFile: ".github/workflows/cd.yml",
		},
		{
			SourceFile: ".tflint.hcl",
			Data:       nil,
			TargetFile: ".tflint.hcl",
		},
	}
}

func TfEnvCICDPreset(tfEnvRepoOpts git.Opts, cicdRepoOpts git.Opts) []Template {
	return []Template{
		{
			SourceFile: fmt.Sprintf("%s/%s_ci.yml", tfEnvRepoOpts.Provider, template.TerraformModule),
			Data: template.CICDTerraformModuleTemplate{
				Project:       cicdRepoOpts.Project,
				Repo:          cicdRepoOpts.Repo,
				Modules:       []string{"base", "k8s"},
				DefaultBranch: cicdRepoOpts.GetDefaultBranch(),
			},
			TargetFile: ".github/workflows/ci.yml",
		},
		{
			SourceFile: fmt.Sprintf("%s/%s_cd.yml", tfEnvRepoOpts.Provider, template.TerraformModule),
			Data: template.CICDTerraformModuleTemplate{
				Project:       cicdRepoOpts.Project,
				Repo:          cicdRepoOpts.Repo,
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
	}
}
