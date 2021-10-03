package git

import (
	"fmt"
	"strings"
)

type RemoteActor interface {
	Commit(content *string, file *string, branch *string, message *string, overwrite bool) error
	RequestReview(branch *string, summary *string) error
}

func New(opts *RemoteOpts) (RemoteActor, error) {
	if strings.Contains(opts.URL, "github.com") {
		return newGitHubActor(opts)
	}
	return nil, fmt.Errorf("unsupported Git provider")
}
