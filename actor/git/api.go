package git

import "fmt"

type Actor interface {
	Commit(content *string, file *string, branch *string, message *string, overwrite bool) error
}

type Opts struct {
	Provider string
	Project  string
	Repo     string

	RemoteBaseURL  string
	RemoteAuthUser string
	RemoteAuthPass string

	DefaultBranch *string
	AuthorName    *string
	AuthorEmail   *string
}

var defaultOpts = Opts{
	DefaultBranch: ptr("main"),
	AuthorName:    ptr("bootstrapper"),
	AuthorEmail:   ptr("bootstrapper@example.com"),
}

func (o *Opts) GetRemoteURL() string {
	return fmt.Sprintf("https://%s/%s/%s", o.Provider, o.Project, o.Repo)
}

func (o *Opts) GetDefaultBranch() string {
	if o.DefaultBranch == nil {
		return *defaultOpts.DefaultBranch
	}
	return *o.DefaultBranch
}

func (o *Opts) GetAuthorName() string {
	if o.AuthorName == nil {
		return *defaultOpts.AuthorName
	}
	return *o.AuthorName
}

func (o *Opts) GetAuthorEmail() string {
	if o.AuthorEmail == nil {
		return *defaultOpts.AuthorEmail
	}
	return *o.AuthorEmail
}

func ptr(v string) *string {
	vv := v
	return &vv
}
