package git

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
	"net/http"
	"sync"
)

type gitHubActor struct {
	Owner  string
	Repo   string
	client *github.Client
	once   sync.Once
}

func (g *gitHubActor) Commit(content *string, file *string, branch *string, message *string, overwrite bool) error {
	if g.client == nil {
		return fmt.Errorf("not authenticated")
	}

	ctx := context.TODO()

	defaultBranch, err := g.getDefaultBranch(ctx)
	if err != nil {
		return err
	}

	oldFileContent, err := g.getFileContent(ctx, file, defaultBranch)
	if err != nil {
		return err
	}

	err = g.createBranch(ctx, branch, defaultBranch)
	if err != nil {
		return err
	}

	if oldFileContent == nil {
		err = g.createOrUpdateFile(ctx, content, file, branch, message, nil)
		if err != nil {
			return err
		}
	} else {
		if overwrite {
			err = g.createOrUpdateFile(ctx, content, file, branch, message, oldFileContent.SHA)
			if err != nil {
				return err
			}
		} else {
			dec, err := base64.StdEncoding.DecodeString(*oldFileContent.Content)
			if err != nil {
				return err
			}
			newContent := fmt.Sprintf("%s\n%s", dec, *content)
			err = g.createOrUpdateFile(ctx, &newContent, file, branch, message, oldFileContent.SHA)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *gitHubActor) getDefaultBranch(ctx context.Context) (*string, error) {
	repo, _, err := g.client.Repositories.Get(ctx, g.Owner, g.Repo)
	if err != nil {
		return nil, err
	}
	return repo.DefaultBranch, nil
}

func (g *gitHubActor) getFileContent(ctx context.Context, file *string, branch *string) (*github.RepositoryContent, error) {
	opts := github.RepositoryContentGetOptions{Ref: *branch}
	fileContent, _, resp, err := g.client.Repositories.GetContents(ctx, g.Owner, g.Repo, *file, &opts)

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return fileContent, nil
}

func (g *gitHubActor) createBranch(ctx context.Context, newBranch *string, defaultBranch *string) error {
	defaultRef, _, err := g.client.Git.GetRef(ctx, g.Owner, g.Repo, "heads/"+*defaultBranch)
	if err != nil {
		return err
	}

	newRef := github.Reference{
		Ref:    github.String("heads/" + *newBranch),
		Object: defaultRef.Object,
	}
	_, _, err = g.client.Git.CreateRef(ctx, g.Owner, g.Repo, &newRef)
	return err
}

func (g *gitHubActor) createOrUpdateFile(ctx context.Context, content *string, file *string, branch *string, message *string, oldSHA *string) error {
	opts := &github.RepositoryContentFileOptions{
		Message:   message,
		Content:   []byte(*content),
		SHA:       oldSHA,
		Branch:    branch,
		Committer: &github.CommitAuthor{Name: github.String("bootstrapper"), Email: github.String("bootstrapper@example.com")},
	}

	_, _, err := g.client.Repositories.CreateFile(ctx, g.Owner, g.Repo, *file, opts)
	return err
}

func (g *gitHubActor) RequestReview(branch *string, summary *string) error {
	if g.client == nil {
		return fmt.Errorf("not authenticated")
	}

	ctx := context.TODO()
	defaultBranch, err := g.getDefaultBranch(ctx)
	if err != nil {
		return err
	}

	newPR := &github.NewPullRequest{
		Title: summary,
		Head:  branch,
		Base:  defaultBranch,
	}

	pr, _, err := g.client.PullRequests.Create(ctx, g.Owner, g.Repo, newPR)
	if err != nil {
		return err
	}
	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}

func (g *gitHubActor) Authenticate(token *string) {
	g.once.Do(func() {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *token},
		)
		tc := oauth2.NewClient(ctx, ts)

		g.client = github.NewClient(tc)
	})
}
