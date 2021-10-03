package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io/ioutil"
	"os"
	"time"
)

type LocalActor interface {
	Actor
	Init(path string) error
}

type localActor struct {
	Opts
	r *git.Repository
}

func NewLocal(opts *Opts) LocalActor {
	return &localActor{Opts: *opts}
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
		Create: true,
		Branch: plumbing.ReferenceName(*branch),
	})
	if err != nil {
		err = l.r.CreateBranch(&config.Branch{
			Name:   *branch,
			Remote: "origin",
			Merge:  plumbing.ReferenceName("refs/heads/" + *branch),
		})
		if err != nil {
			return err
		}
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
