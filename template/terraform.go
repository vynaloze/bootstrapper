package template

import (
	_ "embed"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"regexp"
	"sort"
)

//go:embed templates/gitignore_editors.tpl
var gitignoreEditors string

//go:embed templates/gitignore_terraform.tpl
var gitignoreTerraform string

func TerraformGitignore() string {
	return gitignoreTerraform + "\n" + gitignoreEditors
}

func TerraformModuleCall(name string, source string, vars map[string]interface{}) (string, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	moduleBlock := rootBody.AppendNewBlock("module", []string{name})
	moduleBody := moduleBlock.Body()
	moduleBody.SetAttributeValue("source", cty.StringVal(source))
	moduleBody.AppendNewline()
	for k, v := range vars { // FIXME keys are not in order (well, obviously) and it's not possible to set whitespace
		t, err := gocty.ImpliedType(v)
		if err != nil {
			return "", err
		}
		vv, err := gocty.ToCtyValue(v, t)
		if err != nil {
			return "", err
		}
		moduleBody.SetAttributeValue(k, vv)
	}
	rootBody.AppendNewline()
	return string(f.Bytes()), nil
}

type TerraformProvider struct {
	Source  string
	Version string
}

func TerraformVersionsFile(providers map[string]TerraformProvider) (string, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	terraformBlock := rootBody.AppendNewBlock("terraform", nil)
	terraformBody := terraformBlock.Body()
	providersBlock := terraformBody.AppendNewBlock("required_providers", nil)
	providersBody := providersBlock.Body()
	for _, name := range sortedKeys(providers) {
		// There should be just:
		// providersBody.SetAttributeValue(name, cty.ObjectVal(map[string]cty.Value{...}))
		// but hclwrite is not able to produce a map where every entry is on a new line,
		// instead it puts entire map in a single line. This is a workaround for pretty-formatting:
		// produce block now and later just add "=" between block name and "{" (using blocksAsMap() function).
		providerBlock := providersBody.AppendNewBlock(name, nil)
		providerBody := providerBlock.Body()
		providerBody.SetAttributeValue("source", cty.StringVal(providers[name].Source))
		providerBody.SetAttributeValue("version", cty.StringVal(providers[name].Version))
	}
	return blocksAsMap(string(f.Bytes()), providers)
}

func blocksAsMap(body string, providers map[string]TerraformProvider) (string, error) {
	for name := range providers {
		re, err := regexp.Compile(name + " {")
		if err != nil {
			return "", err
		}
		body = re.ReplaceAllString(body, name+" = {")
	}
	return body, nil
}

func sortedKeys(m map[string]TerraformProvider) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

type TfInfraSharedCoreTfVars struct {
	TfInfraRepos map[string]TfInfraSharedCoreTfVarsRepo `hcl:"tf_infra_repos"`

	TfcOrgName   string `hcl:"tfc_org_name"`
	RepoOwner    string `hcl:"repo_owner"`
	RepoUser     string `hcl:"repo_user"`
	RepoPassword string `hcle:"omit"`
}
type TfInfraSharedCoreTfVarsRepo struct {
	DefaultBranch string `hcl:"default_branch"`
	Strict        bool   `hcl:"strict"`
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
