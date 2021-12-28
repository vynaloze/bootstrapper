package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"log"
	"time"
)

type AddCICDToRepoOpts struct {
	TargetRepoOpts git.Opts
	Templates      []Template
}

func AddCICDToRepo(opts *AddCICDToRepoOpts) error {
	log.Printf("adding CICD to %s repo", opts.TargetRepoOpts.Repo)

	localActor := git.NewLocal(&opts.TargetRepoOpts)
	remoteActor, err := git.NewRemote(&opts.TargetRepoOpts)
	if err != nil {
		return fmt.Errorf("cannot initialize remote Git actor: %w", err)
	}

	log.Printf("preparing CICD pipelines templates")

	gitFiles := make([]git.File, 0)
	for _, file := range opts.Templates {
		filename := fmt.Sprintf("templates/cicd/pipeline_templates/%s", file.SourceFile)
		var pipelineFile []byte
		if file.Data == nil {
			pipelineFile, err = template.Raw(filename)
		} else {
			pipelineFile, err = template.Parse(filename, file.Data)
		}
		if err != nil {
			return fmt.Errorf("error fetching template: %w", err)
		}
		gitFiles = append(gitFiles, git.File{Filename: file.TargetFile, Content: string(pipelineFile)})
	}

	log.Printf("pushing changes to remote repository")

	branch := fmt.Sprintf("%s/%d", opts.TargetRepoOpts.GetAuthorName(), time.Now().UnixMilli())
	message := "chore: add CI/CD pipelines templates"

	err = localActor.CommitMany(branch, message, gitFiles...)
	if err != nil {
		return fmt.Errorf("error committing files: %w", err)
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
