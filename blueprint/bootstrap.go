package blueprint

import (
	"bootstrapper/template"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const defaultBranch = "main"

type BootstrapOpts struct {
}

func Bootstrap() error {
	_, err := initLocalRepo()
	if err != nil {
		return err
	}
	return nil
}

func initLocalRepo() (*git.Worktree, error) {
	dir, err := os.MkdirTemp("", "bootstrapper-")
	if err != nil {
		return nil, err
	}

	repoDir := filepath.Join(dir, "TerraformInfraSharedRepoName") // FIXME
	r, err := git.PlainInit(repoDir, false)
	if err != nil {
		return nil, err
	}
	fmt.Println("repo path: " + repoDir) // TODO

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	err = r.CreateBranch(&config.Branch{
		Name:   defaultBranch,
		Remote: "origin",
		Merge:  "refs/heads/" + defaultBranch,
	})
	if err != nil {
		return nil, err
	}

	file := ".gitignore"
	filename := filepath.Join(repoDir, file)
	err = ioutil.WriteFile(filename, []byte(template.TerraformGitignore()), 0644)
	if err != nil {
		return nil, err
	}

	_, err = w.Add(file)
	if err != nil {
		return nil, err
	}

	_, err = w.Commit("add "+file, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "bootstrapper",
			Email: "bootstrapper@example.com",
			When:  time.Now(),
		},
	})

	return w, err
}
