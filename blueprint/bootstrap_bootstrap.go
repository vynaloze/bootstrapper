package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/actor/terraform"
	"bootstrapper/template"
	"fmt"
	hclencoder_blocks "github.com/rodaine/hclencoder"
	hclencoder_maps "github.com/vdombrovski/hclencoder"
	"log"
	"os"
	"path/filepath"
)

type BootstrapOpts struct {
	SharedInfraRepoOpts git.Opts
	CICDRepoOpts        git.Opts
	TerraformOpts

	localRepoDir *string
}

func (b *BootstrapOpts) getLocalRepoDir() (string, error) {
	if b.localRepoDir == nil {
		dir, err := os.MkdirTemp("", b.SharedInfraRepoOpts.GetAuthorName()+"-")
		if err != nil {
			return "", err
		}
		repoDir := filepath.Join(dir, b.SharedInfraRepoOpts.Repo)
		b.localRepoDir = &repoDir
	}
	return *b.localRepoDir, nil
}

func Bootstrap(opts *BootstrapOpts) error {
	log.Printf("starting bootstrap process")
	gitActor := git.NewLocal(&opts.SharedInfraRepoOpts)
	tfActor, err := terraform.New(&opts.Opts)
	if err != nil {
		return fmt.Errorf("cannot initialize Terraform binary: %w", err)
	}

	log.Printf("initializing local repository: %s", opts.SharedInfraRepoOpts.Repo)
	err = initLocalRepo(gitActor, opts)
	if err != nil {
		return err
	}
	log.Printf("rendering terraform code")
	err = renderTerraformCode(gitActor, opts)
	if err != nil {
		return err
	}
	log.Printf("executing terraform apply")
	err = localApply(*tfActor, opts)
	if err != nil {
		return err
	}
	log.Printf("pushing changes to created remote repository")
	err = gitActor.Push()
	if err != nil {
		return err
	}

	return nil
}

func initLocalRepo(gitActor git.LocalActor, opts *BootstrapOpts) error {
	repoDir, err := opts.getLocalRepoDir()
	if err != nil {
		return err
	}

	err = gitActor.Init(repoDir)
	if err != nil {
		return err
	}
	log.Printf("local repo path: %s", repoDir)

	file := ".gitignore"
	content := template.TerraformGitignore()
	branch := opts.SharedInfraRepoOpts.GetDefaultBranch()
	message := "add .gitignore"

	return gitActor.Commit(&content, &file, &branch, &message, false)
}

func renderTerraformCode(gitActor git.LocalActor, opts *BootstrapOpts) error {
	terraformTfContent, err := renderTerraformTfContent(opts)
	if err != nil {
		return err
	}
	tfVarsContent, err := renderTfVarsContent(opts)
	if err != nil {
		return err
	}
	reposTfContent, err := renderRaw("repos.tf", opts)
	if err != nil {
		return err
	}
	tfCloudTfContent, err := renderRaw("tf_cloud.tf", opts)
	if err != nil {
		return err
	}
	variablesTfContent, err := renderRaw("variables.tf", opts)
	if err != nil {
		return err
	}
	versionsTfContent, err := renderRaw("versions.tf", opts)
	if err != nil {
		return err
	}

	files := []git.File{
		{filepath.Join(opts.GetTerraformInfraCoreDir(), "repos.tf"), string(reposTfContent)},
		{filepath.Join(opts.GetTerraformInfraCoreDir(), "terraform.tf"), string(terraformTfContent)},
		{filepath.Join(opts.GetTerraformInfraCoreDir(), "tf_cloud.tf"), string(tfCloudTfContent)},
		{filepath.Join(opts.GetTerraformInfraCoreDir(), "variables.tf"), string(variablesTfContent)},
		{filepath.Join(opts.GetTerraformInfraCoreDir(), "versions.tf"), string(versionsTfContent)},
		{filepath.Join(opts.GetTerraformInfraCoreDir(), "terraform.auto.tfvars"), string(tfVarsContent)},
	}
	branch := opts.SharedInfraRepoOpts.GetDefaultBranch()
	message := fmt.Sprintf("feat: add %s and %s repos", opts.SharedInfraRepoOpts.Repo, opts.CICDRepoOpts.Repo)

	return gitActor.CommitMany(branch, message, files...)
}

func renderRaw(file string, opts *BootstrapOpts) ([]byte, error) {
	return template.Raw(fmt.Sprintf("templates/tf-infra-shared/%s/core/%s", opts.SharedInfraRepoOpts.Provider, file))
}

func renderTerraformTfContent(opts *BootstrapOpts) ([]byte, error) {
	terraformTf := template.TfInfraSharedCoreTerraformTf{
		Terraform: template.TfInfraSharedCoreTerraformTfTerraform{
			Backend: template.TfInfraSharedCoreTerraformTfBackend{
				Name:         "remote",
				Hostname:     "app.terraform.io",
				Organization: opts.TerraformCloudOrg,
				Workspaces: template.TfInfraSharedCoreTerraformTfWorkspaces{
					Name: fmt.Sprintf("%s-%s", opts.SharedInfraRepoOpts.Repo, opts.GetTerraformInfraCoreDir()),
				},
			},
		},
	}
	return hclencoder_blocks.Encode(terraformTf)
}

func renderTfVarsContent(opts *BootstrapOpts) ([]byte, error) {
	tfVars := template.TfInfraSharedCoreTfVars{
		TfInfraRepos: map[string]template.TfInfraSharedCoreTfVarsRepo{
			opts.SharedInfraRepoOpts.Repo: {
				Modules:       []string{opts.GetTerraformInfraCoreDir()},
				DefaultBranch: opts.SharedInfraRepoOpts.GetDefaultBranch(),
				Strict:        true,
				BuildChecks:   []string{"terraform / ci"},
			},
		},
		TfModuleRepos: map[string]template.TfInfraSharedCoreTfVarsRepo{},
		MiscRepos: map[string]template.TfInfraSharedCoreTfVarsRepo{
			opts.CICDRepoOpts.Repo: {
				Modules:       []string{},
				DefaultBranch: opts.CICDRepoOpts.GetDefaultBranch(),
				Strict:        true,
				BuildChecks:   []string{},
			},
		},
		TfcOrganization: opts.TerraformCloudOrg,
	}
	return hclencoder_maps.Encode(tfVars)
}

func localApply(tfActor terraform.Actor, opts *BootstrapOpts) error {
	repoDir, err := opts.getLocalRepoDir()
	if err != nil {
		return err
	}
	return tfActor.Apply(filepath.Join(repoDir, opts.GetTerraformInfraCoreDir()))
}
