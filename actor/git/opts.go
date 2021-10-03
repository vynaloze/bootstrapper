package git

type Opts struct {
	DefaultBranch *string
	AuthorName    *string
	AuthorEmail   *string
}

var defaultOpts = Opts{
	DefaultBranch: ptr("main"),
	AuthorName:    ptr("bootstrapper"),
	AuthorEmail:   ptr("bootstrapper@example.com"),
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

type RemoteOpts struct {
	Opts
	URL  string
	Auth string
}

func ptr(v string) *string {
	vv := v
	return &vv
}
