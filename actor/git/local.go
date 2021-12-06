package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type LocalActor interface {
	Actor
	Init(path string) error
	Push(repoName string) error
}

type localActor struct {
	Opts
	r *git.Repository
}

func NewLocal(opts Opts) LocalActor {
	return &localActor{Opts: opts}
}

func (l *localActor) Commit(content *string, file *string, branch *string, message *string, overwrite bool) error {
	if _, err := os.Stat(*file); !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("file %s already exists and overwrite=false", *file)
	}

	w, err := l.r.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Create: false,
		Branch: plumbing.NewBranchReferenceName(*branch),
	})
	if err != nil {
		if _, err := l.r.Head(); err != nil {
			// before the first commit - skip
		} else {
			err = w.Checkout(&git.CheckoutOptions{
				Create: true,
				Branch: plumbing.NewBranchReferenceName(*branch),
			})
			if err != nil {
				return err
			}
		}
	}

	parent := filepath.Dir(*file)
	if parent != string(os.PathSeparator) && parent != "." {
		// needs to create a directory
		_ = os.Mkdir(parent, 0644) // ignore errors if directory exists
	}

	err = ioutil.WriteFile(*file, []byte(*content), 0644)
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit(*message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  l.GetAuthorName(),
			Email: l.GetAuthorEmail(),
			When:  time.Now(),
		},
	})
	return err
}

func (l *localActor) Init(path string) error {
	r, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}
	l.r = r
	return nil
}

func (l *localActor) Push(repoName string) error {
	_, err := l.r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{l.RemoteBaseURL + "/" + repoName},
	})
	if err != nil {
		return err
	}
	return l.r.Push(&git.PushOptions{Auth: &http.BasicAuth{Username: l.RemoteAuthUser, Password: l.RemoteAuthPass}})
}
