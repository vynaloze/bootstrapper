package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"strings"
	"time"
)

func CreateApplicationGitRepo(opts ApplicationGitRepoOpts) error {
	if strings.Count(opts.RemoteOpts.URL, "/") < 2 {
		opts.RemoteOpts.URL = opts.RemoteOpts.URL + "/" + opts.GetSharedInfraRepoName()
	}
	gitActor, err := git.New(&opts.RemoteOpts)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{"name": opts.RepoName, "private": true}
	renderedTemplate, err := template.TerraformModuleCall("git_repo_"+opts.RepoName, opts.GetRepoModuleSource(), vars)
	if err != nil {
		return err
	}

	targetFile := opts.GetSharedInfraCoreDir() + "/repos.tf"
	branch := fmt.Sprintf("%s/%s/%d", opts.GetAuthorName(), opts.RepoName, time.Now().UnixMilli())
	message := fmt.Sprintf("[%s] add git repo for %s", opts.GetAuthorName(), opts.RepoName)

	err = gitActor.Commit(&renderedTemplate, &targetFile, &branch, &message, false)
	if err != nil {
		return err
	}
	return gitActor.RequestReview(&branch, &message)
}

func AddApplicationGitRepos() {} //add to existing one (for_each-ed)???

func CreateApplicationsGitRepos() {}

type ApplicationGitRepoOpts struct {
	git.RemoteOpts
	TerraformOpts
	RepoName          string
	RepoModuleSource  *string
	RepoModuleVersion *string
}

var defaultOpts = ApplicationGitRepoOpts{
	RepoModuleSource:  ptr("git::git@%s:%s/terraform-github-repository-manual.git?ref=%s"), // FIXME
	RepoModuleVersion: ptr("v1.0.0"),
}

func (o *ApplicationGitRepoOpts) GetRepoModuleSource() string {
	if o.RepoModuleSource == nil {
		return *defaultOpts.RepoModuleSource
	}
	return *o.RepoModuleSource
}

func (o *ApplicationGitRepoOpts) GetRepoModuleVersion() string {
	if o.RepoModuleVersion == nil {
		return *defaultOpts.RepoModuleVersion
	}
	return *o.RepoModuleVersion
}

func ptr(v string) *string {
	vv := v
	return &vv
}
