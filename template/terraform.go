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

type TfInfraBaseTfVars struct {
	Region              string   `hcl:"region"`
	Environment         string   `hcl:"environment"`
	BaseDomain          string   `hcl:"base_domain"`
	VpcCidr             string   `hcl:"vpc_cidr"`
	ClientCidrBlock     string   `hcl:"client_cidr_block"`
	PrivateSubnetsCidrs []string `hcl:"private_subnets_cidrs"`
	PublicSubnetsCidrs  []string `hcl:"public_subnets_cidrs"`
}

type TfInfraK8sTfVars struct {
	Region              string   `hcl:"region"`
	Environment         string   `hcl:"environment"`
	VpcCidr             string   `hcl:"vpc_cidr"`
	PrivateSubnetsCidrs []string `hcl:"private_subnets_cidrs"`
	PublicSubnetsCidrs  []string `hcl:"public_subnets_cidrs"`
}

type TfInfraDataTf struct {
	RemoteBackends []TfInfraDataTfRemoteBackend `hcl:"data,block"`
}
type TfInfraDataTfRemoteBackend struct {
	Type string `hcl:",key"`
	Name string `hcl:",key"`

	Backend string                           `hcl:"backend"`
	Config  TfInfraDataTfRemoteBackendConfig `hcl:"config"`
}
type TfInfraDataTfRemoteBackendConfig struct {
	Organization string                                     `hcl:"organization"`
	Workspaces   TfInfraDataTfRemoteBackendConfigWorkspaces `hcl:"workspaces"`
}
type TfInfraDataTfRemoteBackendConfigWorkspaces struct {
	Name string `hcl:"name"`
}

type TfInfraTerraformTf struct {
	Terraform TfInfraTerraformTfTerraform `hcl:"terraform,block"`
}
type TfInfraTerraformTfTerraform struct {
	Backend TfInfraTerraformTfBackend `hcl:"backend,block"`
}
type TfInfraTerraformTfBackend struct {
	Name         string                       `hcl:",key"`
	Hostname     string                       `hcl:"hostname"`
	Organization string                       `hcl:"organization"`
	Workspaces   TfInfraTerraformTfWorkspaces `hcl:"workspaces,block"`
}
type TfInfraTerraformTfWorkspaces struct {
	Name string `hcl:"name"`
}
