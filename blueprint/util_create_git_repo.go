package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"github.com/hashicorp/hcl"
	hclencoder_maps "github.com/vdombrovski/hclencoder"
	"log"
	"time"
)

type CreateGitRepoOpts struct {
	SharedInfraRepoOpts git.Opts
	NewRepoOpts         git.Opts
	NewRepoType         template.GitRepoType

	// optional
	NewRepoExtraContent template.GitRepoExtraContent
}

func CreateGitRepo(opts CreateGitRepoOpts) error {
	log.Printf("creating new git repo: %s (type: %s) using %s repo", opts.NewRepoOpts.Repo, opts.NewRepoType, opts.SharedInfraRepoOpts.Repo)
	sharedInfraGitActor, err := git.NewRemote(&opts.SharedInfraRepoOpts)
	if err != nil {
		return fmt.Errorf("error initializing git actor: %w", err)
	}

	log.Printf("reading existing repos")
	file := "core/terraform.auto.tfvars"
	reposContent, err := sharedInfraGitActor.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading existing repos from remote: %w", err)
	}
	var reposTfVars template.TfInfraSharedCoreTfVars
	err = hcl.Decode(&reposTfVars, reposContent)
	if err != nil {
		return fmt.Errorf("error decoding existing repos: %w", err)
	}
	log.Printf("found: %+v", reposTfVars)

	log.Printf("adding new repo")
	err = reposTfVars.AddRepo(opts.NewRepoType, opts.NewRepoExtraContent, opts.NewRepoOpts)
	if err != nil {
		return fmt.Errorf("error adding new repo: %w", err)
	}
	log.Printf("pushing changes to remote repository")

	content, err := hclencoder_maps.Encode(reposTfVars)
	if err != nil {
		return fmt.Errorf("error encoding updated tfvars: %w", err)
	}
	contentStr := string(content)
	branch := fmt.Sprintf("%s/%d", opts.NewRepoOpts.GetAuthorName(), time.Now().UnixMilli())
	message := fmt.Sprintf("feat: add %s repo", opts.NewRepoOpts.Repo)

	err = sharedInfraGitActor.Commit(&contentStr, &file, &branch, &message, true)
	if err != nil {
		return fmt.Errorf("error committing file: %w", err)
	}
	err = sharedInfraGitActor.RequestReview(&branch, &message)
	if err != nil {
		return fmt.Errorf("error creating PR: %w", err)
	}
	return nil
}
