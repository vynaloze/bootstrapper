package git

import (
	"fmt"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type LocalActor interface {
	Actor
	CommitMany(branch string, message string, files ...File) error
	Init(path string) error
	Push() error
}

type File struct {
	Filename string
	Content  string
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

	return l.CommitMany(*branch, *message, File{*file, *content})
}

func (l *localActor) CommitMany(branch string, message string, files ...File) error {
	w, err := l.r.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Create: false,
		Branch: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		if _, err := l.r.Head(); err != nil {
			// before the first commit - skip
		} else {
			err = w.Checkout(&git.CheckoutOptions{
				Create: true,
				Branch: plumbing.NewBranchReferenceName(branch),
			})
			if err != nil {
				return err
			}
		}
	}

	for _, file := range files {
		parent := filepath.Dir(file.Filename)
		if parent != string(os.PathSeparator) && parent != "." {
			// needs to create a directory
			_ = os.Mkdir(parent, 0644) // ignore errors if directory exists
		}

		err = ioutil.WriteFile(file.Filename, []byte(file.Content), 0644)
		if err != nil {
			return err
		}
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit(message, &git.CommitOptions{
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

func (l *localActor) Push() error {
	_, err := l.r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{l.GetRemoteURL()},
	})
	if err != nil {
		return err
	}

	err = l.r.Fetch(&git.FetchOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}

	headRef, err := l.r.Head()
	if err != nil {
		return err
	}
	branch := headRef.Name().Short()

	// go-git does not support rebase so fuck it
	dir := filepath.Dir(l.r.Storer.(*filesystem.Storage).Filesystem().(*chroot.ChrootHelper).Root())
	cmd := exec.Command("git", "branch", "--set-upstream-to=origin/"+branch, branch)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		log.Println(string(out))
		return err
	}

	cmd = exec.Command("git", "pull", "--rebase")
	cmd.Dir = dir
	out, err = cmd.Output()
	if err != nil {
		log.Println(string(out))
		return err
	}

	return l.r.Push(&git.PushOptions{Auth: &http.BasicAuth{Username: l.RemoteAuthUser, Password: l.RemoteAuthPass}})
}
