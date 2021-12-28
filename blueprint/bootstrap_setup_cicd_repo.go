package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"log"
	"time"
)

type SetupCICDRepoOpts struct {
	CICDRepoOpts git.Opts
	// TODO decide what templates to render (or all?)
}

func SetupCICDRepo(opts *SetupCICDRepoOpts) error {
	log.Printf("setting up CICD repo")

	localActor := git.NewLocal(&opts.CICDRepoOpts)
	remoteActor, err := git.NewRemote(&opts.CICDRepoOpts)
	if err != nil {
		return fmt.Errorf("cannot initialize remote Git actor: %w", err)
	}

	log.Printf("preparing CICD pipelines templates")

	allFiles, err := template.RawAll("templates/cicd")
	if err != nil {
		return fmt.Errorf("error fetching templates: %w", err)
	}
	gitFiles := make([]git.File, 0, len(allFiles))
	for filename, content := range allFiles {
		gitFiles = append(gitFiles, git.File{Filename: filename, Content: string(content)})
	}

	log.Printf("pushing changes to remote repository")

	branch := fmt.Sprintf("%s/%d", opts.CICDRepoOpts.GetAuthorName(), time.Now().UnixMilli())
	message := "feat: add CI/CD pipelines templates"

	err = localActor.CommitMany(branch, message, gitFiles...)
	if err != nil {
		return fmt.Errorf("error committing template files: %w", err)
	}
	err = localActor.Push()
	if err != nil {
		return fmt.Errorf("error pushing changes: %w", err)
	}
	err = remoteActor.RequestReview(&branch, &message)
	if err != nil {
		return fmt.Errorf("error creating PR: %w", err)
	}

	return nil
}
