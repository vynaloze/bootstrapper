package git

import (
	"fmt"
)

type RemoteActor interface {
	Actor
	RequestReview(branch *string, summary *string) error
}

func NewRemote(opts *Opts) (RemoteActor, error) {
	if opts.Provider == "github.com" {
		return newGitHubActor(opts)
	}
	return nil, fmt.Errorf("unsupported Git provider: %s", opts.Provider)
}
