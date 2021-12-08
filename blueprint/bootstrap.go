package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/actor/terraform"
	"bootstrapper/template"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type BootstrapOpts struct {
	git.Opts
	TerraformOpts

	localRepoDir *string
}

func (b *BootstrapOpts) getLocalRepoDir() (string, error) {
	if b.localRepoDir == nil {
		dir, err := os.MkdirTemp("", b.GetAuthorName()+"-")
		if err != nil {
			return "", err
		}
		repoDir := filepath.Join(dir, b.GetSharedInfraRepoName())
		b.localRepoDir = &repoDir
	}
	return *b.localRepoDir, nil
}

func Bootstrap(opts *BootstrapOpts) error {
	log.Printf("starting bootstrap process")
	gitActor := git.NewLocal(opts.Opts)
	tfActor, err := terraform.New()
	if err != nil {
		return fmt.Errorf("cannot initialize Terraform binary: %w", err)
	}

	log.Printf("initializing local repository: %s", opts.GetSharedInfraRepoName())
	err = initLocalRepo(gitActor, opts)
	if err != nil {
		return err
	}
	log.Printf("rendering terraform code")
	err = callTerraformRepoModule(gitActor, opts)
	if err != nil {
		return err
	}
	log.Printf("executing terraform apply")
	err = localApply(*tfActor, opts)
	if err != nil {
		return err
	}
	log.Printf("pushing changes to created remote repository")
	err = gitActor.Push(opts.GetSharedInfraRepoName())
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

	file := filepath.Join(repoDir, ".gitignore")
	content := template.TerraformGitignore()
	branch := opts.GetDefaultBranch()
	message := "add .gitignore"

	return gitActor.Commit(&content, &file, &branch, &message, false)
}

func callTerraformRepoModule(gitActor git.LocalActor, opts *BootstrapOpts) error {
	repoDir, err := opts.getLocalRepoDir()
	if err != nil {
		return err
	}

	file := filepath.Join(repoDir, opts.GetTerraformModuleReposFile())
	content, err := template.TfInfraSharedCoreReposTf(template.TfInfraSharedCoreReposTfOpts{
		Repos:         []string{"tf-infra-shared"},
		Strict:        true, //TODO?
		DefaultBranch: opts.GetDefaultBranch(),
	})
	if err != nil {
		return err
	}
	branch := opts.GetDefaultBranch()
	message := "feat: add tf-shared-infra repo"

	return gitActor.Commit(&content, &file, &branch, &message, false)
}

func localApply(tfActor terraform.Actor, opts *BootstrapOpts) error {
	repoDir, err := opts.getLocalRepoDir()
	if err != nil {
		return err
	}

	for k, v := range opts.ProviderSecrets {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return tfActor.Apply(filepath.Dir(filepath.Join(repoDir, opts.GetTerraformInfraReposFile())))
}
