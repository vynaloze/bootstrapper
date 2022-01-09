package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	hclencoder_blocks "github.com/rodaine/hclencoder"
	hclencoder_maps "github.com/vdombrovski/hclencoder"
	"log"
	"time"
)

type SetupCloudInfraOpts struct {
	InfraRepoOpts git.Opts
	EnvRepoOpts   git.Opts

	CloudProvider     template.CloudProvider
	TerraformCloudOrg string

	CICDTemplates          []Template
	TerraformVars          interface{}
	TerraformParentModules []string

	Module string
}

func SetupCloudInfra(opts *SetupCloudInfraOpts) error {
	log.Printf("setting up cloud infra repo")

	infraLocalActor := git.NewLocal(&opts.InfraRepoOpts)
	infraRemoteActor, err := git.NewRemote(&opts.InfraRepoOpts)
	if err != nil {
		return fmt.Errorf("cannot initialize remote Git actor for infra repo: %w", err)
	}
	envRemoteActor, err := git.NewRemote(&opts.EnvRepoOpts)
	if err != nil {
		return fmt.Errorf("cannot initialize remote Git actor for env repo: %w", err)
	}

	branch := fmt.Sprintf("%s/%d", opts.InfraRepoOpts.GetAuthorName(), time.Now().UnixMilli())

	log.Printf("preparing CICD pipelines templates")
	ciFiles, err := templatesToGitFiles("cicd/pipeline_templates", opts.CICDTemplates)
	if err != nil {
		return fmt.Errorf("error preparing CICD templates: %w", err)
	}
	err = infraLocalActor.CommitMany(branch, "chore: add CI/CD pipelines templates", ciFiles...)
	if err != nil {
		return fmt.Errorf("error committing ci files: %w", err)
	}

	log.Printf("preparing terraform files")

	latestTag, err := envRemoteActor.LatestTag()
	if err != nil {
		return fmt.Errorf("cannot fetch latest tag from env repo: %w", err)
	}
	templates := []Template{
		{
			SourceFile: fmt.Sprintf("%s/%s/main.tf", opts.CloudProvider, opts.Module),
			Data: template.TerraformModuleTemplate{
				GitProvider: opts.EnvRepoOpts.Provider,
				GitProject:  opts.EnvRepoOpts.Project,
				GitRepo:     opts.EnvRepoOpts.Repo,
				Ref:         latestTag,
			},
			TargetFile: fmt.Sprintf("%s/main.tf", opts.Module),
		}, {
			SourceFile: fmt.Sprintf("%s/%s/outputs.tf", opts.CloudProvider, opts.Module),
			TargetFile: fmt.Sprintf("%s/outputs.tf", opts.Module),
		}, {
			SourceFile: fmt.Sprintf("%s/%s/variables.tf", opts.CloudProvider, opts.Module),
			TargetFile: fmt.Sprintf("%s/variables.tf", opts.Module),
		}, {
			SourceFile: fmt.Sprintf("%s/%s/versions.tf", opts.CloudProvider, opts.Module),
			TargetFile: fmt.Sprintf("%s/versions.tf", opts.Module),
		},
	}
	tfTemplateFiles, err := templatesToGitFiles("tf-infra", templates)
	if err != nil {
		return fmt.Errorf("error preparing templated terraform files: %w", err)
	}
	tfDynamicFiles, err := tfInfraDynamicFiles(opts)
	if err != nil {
		return fmt.Errorf("error preparing dynamic terraform files: %w", err)
	}
	tfFiles := append(tfTemplateFiles, tfDynamicFiles...)
	tfFiles = append(tfFiles, git.File{Filename: ".gitignore", Content: template.TerraformGitignore()})

	message := fmt.Sprintf("feat: add initial %s %s module", opts.CloudProvider, opts.Module)
	err = infraLocalActor.CommitMany(branch, message, tfFiles...)
	if err != nil {
		return fmt.Errorf("error committing tf files: %w", err)
	}

	log.Printf("pushing changes to remote repository")
	err = infraLocalActor.Push()
	if err != nil {
		return fmt.Errorf("error pushing changes: %w", err)
	}
	err = infraRemoteActor.RequestReview(&branch, &message)
	if err != nil {
		return fmt.Errorf("error creating PR: %w", err)
	}

	return nil
}

func tfInfraDynamicFiles(opts *SetupCloudInfraOpts) ([]git.File, error) {
	gitFiles := make([]git.File, 0)
	terraformTfContent, err := renderInfraTerraformTfContent(opts)
	if err != nil {
		return nil, fmt.Errorf("cannot render terraform.tf: %w", err)
	}
	gitFiles = append(gitFiles, git.File{Filename: fmt.Sprintf("%s/terraform.tf", opts.Module), Content: string(terraformTfContent)})

	tfVarsContent, err := renderInfraTfVarsContent(opts)
	if err != nil {
		return nil, fmt.Errorf("cannot render terraform.auto.tfvars: %w", err)
	}
	gitFiles = append(gitFiles, git.File{Filename: fmt.Sprintf("%s/terraform.auto.tfvars", opts.Module), Content: string(tfVarsContent)})

	if len(opts.TerraformParentModules) > 0 {
		dataTfContent, err := renderInfraDataTfContent(opts)
		if err != nil {
			return nil, fmt.Errorf("cannot render data.tf: %w", err)
		}
		gitFiles = append(gitFiles, git.File{Filename: fmt.Sprintf("%s/data.tf", opts.Module), Content: string(dataTfContent)})
	}

	return gitFiles, nil
}

func renderInfraTerraformTfContent(opts *SetupCloudInfraOpts) ([]byte, error) {
	terraformTf := template.TfInfraTerraformTf{
		Terraform: template.TfInfraTerraformTfTerraform{
			Backend: template.TfInfraTerraformTfBackend{
				Name:         "remote",
				Hostname:     "app.terraform.io",
				Organization: opts.TerraformCloudOrg,
				Workspaces: template.TfInfraTerraformTfWorkspaces{
					Name: fmt.Sprintf("%s-%s", opts.InfraRepoOpts.Repo, opts.Module),
				},
			},
		},
	}
	return hclencoder_blocks.Encode(terraformTf)
}

func renderInfraTfVarsContent(opts *SetupCloudInfraOpts) ([]byte, error) {
	return hclencoder_maps.Encode(opts.TerraformVars)
}

func renderInfraDataTfContent(opts *SetupCloudInfraOpts) ([]byte, error) {
	remotes := make([]template.TfInfraDataTfRemoteBackend, 0)

	for _, module := range opts.TerraformParentModules {
		remotes = append(remotes, template.TfInfraDataTfRemoteBackend{
			Type:    "terraform_remote_state",
			Name:    module,
			Backend: "remote",
			Config: template.TfInfraDataTfRemoteBackendConfig{
				Organization: opts.TerraformCloudOrg,
				Workspaces: template.TfInfraDataTfRemoteBackendConfigWorkspaces{
					Name: fmt.Sprintf("%s-%s", opts.InfraRepoOpts.Repo, module),
				},
			},
		})
	}

	data := template.TfInfraDataTf{RemoteBackends: remotes}
	return hclencoder_maps.Encode(data)
}
