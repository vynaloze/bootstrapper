package git

import (
	"bootstrapper/datasource"
	"fmt"
	"strings"
	"sync"
)

type GitActor interface {
	Commit(content *string, file *string, branch *string, message *string, overwrite bool) error
	RequestReview(branch *string, summary *string) error
}

func New(URL string) (GitActor, error) {
	if strings.Contains(URL, "github.com") {
		token, ok := datasource.Find("actor.git.github.token")
		if !ok {
			return nil, fmt.Errorf("required key not found: actor.git.github.token")
		}
		s := strings.Split(URL, "/")
		g := &gitHubActor{s[1], s[2], nil, sync.Once{}}
		g.Authenticate(&token)
		return g, nil
	}
	return nil, fmt.Errorf("unsupported Git provider")
}
