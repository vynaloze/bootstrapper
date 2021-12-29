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

	// optional: if set, only push to given branch and don't create PR
	TargetRepoBranch *string
}

func AddCICDToRepo(opts *AddCICDToRepoOpts) error {
	log.Printf("adding CICD to %s repo", opts.TargetRepoOpts.Repo)

	localActor := git.NewLocal(&opts.TargetRepoOpts)
	remoteActor, err := git.NewRemote(&opts.TargetRepoOpts)
	if err != nil {
		return fmt.Errorf("cannot initialize remote Git actor: %w", err)
	}

	log.Printf("preparing CICD pipelines templates")

	gitFiles := make([]git.File, 0)
	for _, file := range opts.Templates {
		filename := fmt.Sprintf("templates/cicd/pipeline_templates/%s", file.SourceFile)
		var pipelineFile []byte
		if file.Data == nil {
			pipelineFile, err = template.Raw(filename)
		} else {
			pipelineFile, err = template.Parse(filename, file.Data)
		}
		if err != nil {
			return fmt.Errorf("error fetching template: %w", err)
		}
		gitFiles = append(gitFiles, git.File{Filename: file.TargetFile, Content: string(pipelineFile)})
	}

	log.Printf("pushing changes to remote repository")

	message := "chore: add CI/CD pipelines templates"
	if opts.TargetRepoBranch != nil {
		err = commitAndPush(localActor, *opts.TargetRepoBranch, message, gitFiles)
		if err != nil {
			return err
		}
	} else {
		branch := fmt.Sprintf("%s/%d", opts.TargetRepoOpts.GetAuthorName(), time.Now().UnixMilli())
		err = commitAndPush(localActor, branch, message, gitFiles)
		if err != nil {
			return err
		}
		err = remoteActor.RequestReview(&branch, &message)
		if err != nil {
			return fmt.Errorf("error creating PR: %w", err)
		}
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

func TfInfraSharedCICDPreset(sharedInfraGitOpts git.Opts, cicdRepoOpts git.Opts, terraformOpts TerraformOpts) AddCICDToRepoOpts {
	return AddCICDToRepoOpts{
		TargetRepoOpts: sharedInfraGitOpts,
		Templates: []Template{
			{
				SourceFile: fmt.Sprintf("%s/%s_ci.yml", sharedInfraGitOpts.Provider, template.TerraformInfra),
				Data: template.TerraformInfraTemplate{
					Project:       cicdRepoOpts.Project,
					Repo:          cicdRepoOpts.Repo,
					Modules:       []template.TerraformInfraModuleTemplate{{Name: terraformOpts.GetTerraformInfraCoreDir()}},
					DefaultBranch: cicdRepoOpts.GetDefaultBranch(),
				},
				TargetFile: fmt.Sprintf(".github/workflows/%s.ci.yml", terraformOpts.GetTerraformInfraCoreDir()),
			},
			{
				SourceFile: fmt.Sprintf("%s/%s_cd.yml", sharedInfraGitOpts.Provider, template.TerraformInfra),
				Data: template.TerraformInfraTemplate{
					Project:       cicdRepoOpts.Project,
					Repo:          cicdRepoOpts.Repo,
					Modules:       []template.TerraformInfraModuleTemplate{{Name: terraformOpts.GetTerraformInfraCoreDir()}},
					DefaultBranch: cicdRepoOpts.GetDefaultBranch(),
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
}

func TfEnvCICDPreset(tfEnvRepoOpts git.Opts, tfEnvBranch string, cicdRepoOpts git.Opts) AddCICDToRepoOpts {
	return AddCICDToRepoOpts{
		TargetRepoOpts:   tfEnvRepoOpts,
		TargetRepoBranch: &tfEnvBranch,
		Templates: []Template{
			{
				SourceFile: fmt.Sprintf("%s/%s_ci.yml", tfEnvRepoOpts.Provider, template.TerraformModule),
				Data: template.TerraformModuleTemplate{
					Project:       cicdRepoOpts.Project,
					Repo:          cicdRepoOpts.Repo,
					Modules:       []string{"base", "k8s"},
					DefaultBranch: cicdRepoOpts.GetDefaultBranch(),
				},
				TargetFile: ".github/workflows/ci.yml",
			},
			{
				SourceFile: fmt.Sprintf("%s/%s_cd.yml", tfEnvRepoOpts.Provider, template.TerraformModule),
				Data: template.TerraformModuleTemplate{
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
		},
	}
}
