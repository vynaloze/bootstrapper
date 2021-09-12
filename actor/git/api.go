package git

import (
	"fmt"
	"strings"
	"sync"
)

type GitActor interface {
	Commit(content *string, file *string, branch *string, message *string, overwrite bool) error
	RequestReview(branch *string, summary *string) error
}

func New(URL string, username string, password string) (GitActor, error) {
	if strings.Contains(URL, "github.com") {
		s := strings.Split(URL, "/")
		g := &gitHubActor{s[1], s[2], nil, sync.Once{}}
		g.Authenticate(&password)
		return g, nil
	}
	return nil, fmt.Errorf("unsupported Git provider")
}
