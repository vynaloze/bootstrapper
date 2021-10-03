package git

import (
	"fmt"
	"strings"
)

type RemoteActor interface {
	Actor
	RequestReview(branch *string, summary *string) error
}

type RemoteOpts struct {
	Opts
	URL  string
	Auth string
}

func NewRemote(opts *RemoteOpts) (RemoteActor, error) {
	if strings.Contains(opts.URL, "github.com") {
		return newGitHubActor(opts)
	}
	return nil, fmt.Errorf("unsupported Git provider")
}
