package git

import (
	"fmt"
	"strings"
)

type RemoteActor interface {
	Actor
	RequestReview(branch *string, summary *string) error
}

func NewRemote(opts *Opts) (RemoteActor, error) {
	if strings.Contains(opts.RemoteBaseURL, "github.com") {
		return newGitHubActor(opts)
	}
	return nil, fmt.Errorf("unsupported Git provider")
}
