package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/datasource"
	"bootstrapper/template"
	"fmt"
	"time"
)

const (
	defaultRepoModuleSourceTemplate = "git::git@%s:%s/terraform-github-repository-manual.git?ref=%s" //FIXME
	defaultRepoModuleVersion        = "v1.0.0"
)

func CreateApplicationGitRepo(name string) error {
	gitProvider, ok := datasource.Find("git.provider")
	if !ok {
		return fmt.Errorf("required key not found: git.provider")
	}
	gitProject, ok := datasource.Find("git.project")
	if !ok {
		return fmt.Errorf("required key not found: git.project")
	}

	repoModuleSource := prepareRepoModuleSource(gitProvider, gitProject)

	gitActor, err := prepareGitActor(gitProvider, gitProject)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{"name": name, "private": true}
	renderedTemplate, err := template.TerraformModuleCall("git_repo_"+name, repoModuleSource, vars)
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

func prepareRepoModuleSource(gitProvider, gitProject string) string {
	repoModuleSource, ok := datasource.Find("blueprints.git.application.repo_module.source")
	if !ok {
		repoModuleVersion, ok := datasource.Find("blueprints.git.application.repo_module.version")
		if !ok {
			repoModuleVersion = defaultRepoModuleVersion
		}
		repoModuleSource = fmt.Sprintf(defaultRepoModuleSourceTemplate, gitProvider, gitProject, repoModuleVersion)
	}
	return repoModuleSource
}

func prepareGitActor(gitProvider, gitProject string) (git.GitActor, error) {
	repoName, ok := datasource.Find("terraform.infra_shared_repo_name")
	if !ok {
		repoName = TerraformInfraSharedRepoName
	}
	sharedInfraRepoURL := fmt.Sprintf("%s/%s/%s", gitProvider, gitProject, repoName)

	return git.New(sharedInfraRepoURL)
}
