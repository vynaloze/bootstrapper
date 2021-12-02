package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/actor/terraform"
	"bootstrapper/template"
	"fmt"
	"os"
	"path/filepath"
)

type BootstrapOpts struct {
	git.Opts
	TerraformOpts
	EnvVars map[string]string

	localRepoDir string
}

func Bootstrap(opts BootstrapOpts) error {
	gitActor := git.NewLocal(&opts.Opts)
	tfActor, err := terraform.New()
	if err != nil {
		return err
	}

	err = findLocalRepoDir(&opts)
	if err != nil {
		return err
	}

	err = initLocalRepo(gitActor, opts)
	if err != nil {
		return err
	}
	err = callTerraformRepoModule(gitActor, opts)
	if err != nil {
		return err
	}
	err = localApply(*tfActor, opts)
	if err != nil {
		return err
	}

	return nil
}

func findLocalRepoDir(opts *BootstrapOpts) error {
	dir, err := os.MkdirTemp("", opts.GetAuthorName()+"-")
	if err != nil {
		return err
	}
	opts.localRepoDir = filepath.Join(dir, opts.GetSharedInfraRepoName())
	return nil
}

func initLocalRepo(gitActor git.LocalActor, opts BootstrapOpts) error {
	err := gitActor.Init(opts.localRepoDir)
	if err != nil {
		return err
	}
	fmt.Println("repo path: " + opts.localRepoDir) // TODO

	file := filepath.Join(opts.localRepoDir, ".gitignore")
	content := template.TerraformGitignore()
	branch := opts.GetDefaultBranch()
	message := "add .gitignore"

	return gitActor.Commit(&content, &file, &branch, &message, false)
}

func callTerraformRepoModule(gitActor git.LocalActor, opts BootstrapOpts) error {
	file := filepath.Join(opts.localRepoDir, "core", "repos.tf")
	content, err := template.TfInfraSharedCoreReposTf(template.TfInfraSharedCoreReposTfOpts{
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

func localApply(tfActor terraform.Actor, opts BootstrapOpts) error {
	for k, v := range opts.EnvVars {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return tfActor.Apply(filepath.Join(opts.localRepoDir, "core")) //TODO?
}
