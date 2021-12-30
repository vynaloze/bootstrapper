package template

import (
	"bootstrapper/actor/git"
	_ "embed"
	"fmt"
)

//go:embed templates/gitignore_editors.tpl
var gitignoreEditors string

//go:embed templates/gitignore_terraform.tpl
var gitignoreTerraform string

func TerraformGitignore() string {
	return gitignoreTerraform + "\n" + gitignoreEditors
}

type TerraformProvider struct {
	Source  string
	Version string
}

type TfInfraSharedCoreTfVars struct {
	TfInfraRepos  map[string]TfInfraSharedCoreTfVarsRepo `hcl:"tf_infra_repos"`
	TfModuleRepos map[string]TfInfraSharedCoreTfVarsRepo `hcl:"tf_module_repos"`
	MiscRepos     map[string]TfInfraSharedCoreTfVarsRepo `hcl:"misc_repos"`

	TfcOrganization string `hcl:"tfc_organization"`
}

func (t *TfInfraSharedCoreTfVars) AddRepo(typ GitRepoType, extraContent GitRepoExtraContent, gitOpts git.Opts) error {
	name := gitOpts.Repo
	defaultBranch := gitOpts.GetDefaultBranch()
	switch typ {
	case TerraformInfra:
		t.TfInfraRepos[name] = TfInfraSharedCoreTfVarsRepo{extraContent.Modules, defaultBranch, true, []string{"terraform / ci"}}
	case TerraformModule:
		t.TfModuleRepos[name] = TfInfraSharedCoreTfVarsRepo{[]string{}, defaultBranch, true, []string{"terraform / ci"}}
	case Miscellaneous:
		t.MiscRepos[name] = TfInfraSharedCoreTfVarsRepo{[]string{}, defaultBranch, true, []string{}}
	default:
		return fmt.Errorf("unknown type: " + string(typ))
	}
	return nil
}

type TfInfraSharedCoreTfVarsRepo struct {
	Modules []string `hcl:"modules,omitempty"`

	DefaultBranch string   `hcl:"default_branch"`
	Strict        bool     `hcl:"strict"`
	BuildChecks   []string `hcl:"build_checks"`
}

type TfInfraSharedCoreTerraformTf struct {
	Terraform TfInfraSharedCoreTerraformTfTerraform `hcl:"terraform,block"`
}
type TfInfraSharedCoreTerraformTfTerraform struct {
	Backend TfInfraSharedCoreTerraformTfBackend `hcl:"backend,block"`
}
type TfInfraSharedCoreTerraformTfBackend struct {
	Name         string                                 `hcl:",key"`
	Hostname     string                                 `hcl:"hostname"`
	Organization string                                 `hcl:"organization"`
	Workspaces   TfInfraSharedCoreTerraformTfWorkspaces `hcl:"workspaces,block"`
}
type TfInfraSharedCoreTerraformTfWorkspaces struct {
	Name string `hcl:"name"`
}
