package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"time"
)

func CreateApplicationGitRepo(
	name string,
	sourceModuleRepoURL string, // from datasource
	gitActor git.GitActor, // from service discovery
) error {

	vars := map[string]interface{}{"name": name, "private": true}
	renderedTemplate, err := template.TerraformModuleCall("git_repo_"+name, sourceModuleRepoURL, vars)
	if err != nil {
		return err
	}

	targetFile := "apps/git_repos.tf"
	branch := fmt.Sprintf("%s/%s/%d", "bootstrapper", name, time.Now().UnixMilli())
	message := "[bootstrapper] add git repo for " + name

	err = gitActor.Commit(&renderedTemplate, &targetFile, &branch, &message, false)
	if err != nil {
		return err
	}
	return gitActor.RequestReview(&branch, &message)
}

func AddApplicationGitRepos() {} //add to existing one (for_each-ed)???

func CreateApplicationsGitRepos() {}
