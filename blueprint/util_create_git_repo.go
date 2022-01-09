package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"github.com/hashicorp/hcl"
	hclencoder_maps "github.com/vdombrovski/hclencoder"
	"log"
	"strings"
	"time"
)

type CreateGitReposOpts struct {
	SharedInfraRepoOpts git.Opts
	NewReposSpecs       []CreateGitReposNewRepoSpec

	// optional
	NewRepoExtraContent template.GitRepoExtraContent
}

type CreateGitReposNewRepoSpec struct {
	NewRepoOpts git.Opts
	NewRepoType template.GitRepoType
}

func CreateGitRepos(opts CreateGitReposOpts) error {
	log.Printf("creating new git repo(s) using %s repo:", opts.SharedInfraRepoOpts.Repo)
	reposNames := make([]string, 0)
	for _, r := range opts.NewReposSpecs {
		log.Printf("%s (type: %s)", r.NewRepoOpts.Repo, r.NewRepoType)
		reposNames = append(reposNames, r.NewRepoOpts.Repo)
	}
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

	log.Printf("adding new repo(s)")
	for _, r := range opts.NewReposSpecs {
		err = reposTfVars.AddRepo(r.NewRepoType, opts.NewRepoExtraContent, r.NewRepoOpts)
		if err != nil {
			return fmt.Errorf("error adding new repo: %w", err)
		}
	}
	log.Printf("pushing changes to remote repository")

	content, err := hclencoder_maps.Encode(reposTfVars)
	if err != nil {
		return fmt.Errorf("error encoding updated tfvars: %w", err)
	}
	contentStr := string(content)
	branch := fmt.Sprintf("%s/%d", opts.NewReposSpecs[0].NewRepoOpts.GetAuthorName(), time.Now().UnixMilli())
	message := fmt.Sprintf("feat: add repo(s): %s", strings.Join(reposNames, ", "))

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
