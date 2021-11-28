package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"strings"
	"time"
)

func CreateApplicationGitRepo(opts ApplicationGitRepoOpts) error {
	if strings.Count(opts.RemoteOpts.URL, "/") < 2 {
		opts.RemoteOpts.URL = opts.RemoteOpts.URL + "/" + opts.GetSharedInfraRepoName()
	}
	gitActor, err := git.NewRemote(&opts.RemoteOpts)
	if err != nil {
		return err
	}

	// TODO for_each-ed module
	renderedTemplate, err := template.TerraformModuleFromRegistry(
		"git_repo_"+opts.RepoName,
		opts.RepoModuleSource,
		opts.RepoModuleVersion,
		opts.RepoModuleVars,
	)
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
	RepoModuleSource  string
	RepoModuleVersion string
	RepoModuleVars    []*template.TerraformVariable
}

func (o *ApplicationGitRepoOpts) WithGitHubDefaults(strict bool) {
	o.RepoModuleSource = "mineiros-io/repository/github"
	o.RepoModuleVersion = "~> 0.10.0"
	o.RepoModuleVars = []*template.TerraformVariable{
		{"name", cty.StringVal(o.RepoName)},
		{"visibility", cty.StringVal("private")},
		{"auto_init", cty.BoolVal(false)},
		nil,
		// TODO access control
		// TODO CI integration (Actions are not yet in provider unfortunately)
		{"allow_rebase_merge", cty.BoolVal(true)},
		{"allow_merge_commit", cty.BoolVal(!strict)},
		{"delete_branch_on_merge", cty.BoolVal(true)},
		nil,
		{"branch_protections_v3", cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"branch":         cty.StringVal(o.GetDefaultBranch()),
				"enforce_admins": cty.BoolVal(true),
				"required_status_checks": cty.ObjectVal(map[string]cty.Value{
					"strict": cty.BoolVal(strict),
					"contexts": cty.ListVal([]cty.Value{
						cty.StringVal("ci/TODO"), //TODO requires CI integration
					}),
				}),
				"required_pull_request_reviews": cty.ObjectVal(map[string]cty.Value{
					"dismiss_stale_reviews":           cty.BoolVal(true),
					"require_code_owner_reviews":      cty.BoolVal(true),
					"required_approving_review_count": cty.NumberIntVal(1),
				}),
				//"restrictions": map[string][]string{
				// TODO requires access management
				//},
			}),
		})},
	}
}
