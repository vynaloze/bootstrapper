package git

import (
	"errors"
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
	Clone() (string, error)
}

type File struct {
	Filename string
	Content  string
}

type localActor struct {
	*Opts
	r       *git.Repository
	repoDir string
}

func NewLocal(opts *Opts) LocalActor {
	return &localActor{Opts: opts}
}

func (l *localActor) Commit(content *string, file *string, branch *string, message *string, overwrite bool) error {
	fullPath := filepath.Join(l.repoDir, *file)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("file %s already exists and overwrite=false", fullPath)
	}

	return l.CommitMany(*branch, *message, File{*file, *content})
}

func (l *localActor) CommitMany(branch string, message string, files ...File) error {
	if l.r == nil {
		_, err := l.Clone()
		if err != nil {
			return fmt.Errorf("error cloning repository: %w", err)
		}
	}

	w, err := l.r.Worktree()
	if err != nil {
		return fmt.Errorf("error creating worktree: %w", err)
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
				return fmt.Errorf("error checking out %s: %w", branch, err)
			}
		}
	}

	for _, file := range files {
		fullPath := filepath.Join(l.repoDir, file.Filename)
		parent := filepath.Dir(fullPath)
		if parent != string(os.PathSeparator) && parent != "." {
			// needs to create a directory
			err := os.MkdirAll(parent, 0755)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %w", parent, err)
			}
		}

		err := ioutil.WriteFile(fullPath, []byte(file.Content), 0755)
		if err != nil {
			return fmt.Errorf("error writing file %s: %w", file.Filename, err)
		}
	}

	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("error adding changes to staging area: %w", err)
	}

	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  l.GetAuthorName(),
			Email: l.GetAuthorEmail(),
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("error comitting changes: %w", err)
	}
	return nil
}

func (l *localActor) Clone() (string, error) {
	dir, err := os.MkdirTemp("", l.GetAuthorName()+"-")
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, l.Repo)
	l.repoDir = path

	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:  l.GetRemoteURL(),
		Auth: &http.BasicAuth{Username: l.RemoteAuthUser, Password: l.RemoteAuthPass},
	})
	if err != nil {
		return "", err
	}
	log.Printf("cloned %s to %s", l.GetRemoteURL(), path)
	l.r = r
	return path, nil
}

func (l *localActor) Init(path string) error {
	r, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}
	l.r = r
	l.repoDir = path
	return nil
}

func (l *localActor) Push() error {
	remotes, err := l.r.Remotes()
	if err != nil {
		return err
	}
	if len(remotes) == 0 {
		_, err := l.r.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{l.GetRemoteURL()},
		})
		if err != nil {
			return fmt.Errorf("error creating remote: %w", err)
		}
	}

	err = l.r.Fetch(&git.FetchOptions{RemoteName: "origin"})
	if err != nil {
		if !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return fmt.Errorf("error fetching changes: %w", err)
		}
	} else {
		err := l.rebase()
		if err != nil {
			return fmt.Errorf("error pulling with rebase: %w", err)
		}
	}

	err = l.r.Push(&git.PushOptions{Auth: &http.BasicAuth{Username: l.RemoteAuthUser, Password: l.RemoteAuthPass}})
	if err != nil {
		return fmt.Errorf("error pushing changes: %w", err)
	}
	return nil
}

func (l *localActor) rebase() error {
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
	return nil
}
