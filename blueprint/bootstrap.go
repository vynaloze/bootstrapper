package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"os"
	"path/filepath"
)

type BootstrapOpts struct {
	git.Opts
	TerraformOpts
}

func Bootstrap(opts BootstrapOpts) error {
	gitActor := git.NewLocal(&opts.Opts)

	err := initLocalRepo(gitActor, opts)
	if err != nil {
		return err
	}
	return nil
}

func initLocalRepo(gitActor git.LocalActor, opts BootstrapOpts) error {
	dir, err := os.MkdirTemp("", opts.GetAuthorName()+"-")
	if err != nil {
		return err
	}
	repoDir := filepath.Join(dir, opts.GetSharedInfraRepoName())

	err = gitActor.Init(repoDir)
	if err != nil {
		return err
	}
	fmt.Println("repo path: " + repoDir) // TODO

	file := filepath.Join(repoDir, ".gitignore")
	content := template.TerraformGitignore()
	branch := opts.GetDefaultBranch()
	message := "add .gitignore"

	return gitActor.Commit(&content, &file, &branch, &message, false)
}
