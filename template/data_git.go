package template

import _ "embed"

type GitRepoType string

const (
	TerraformInfra  GitRepoType = "tf_infra"
	TerraformModule GitRepoType = "tf_module"
	Miscellaneous   GitRepoType = "misc"
)

type GitRepoExtraContent struct {
	Modules []string
}

//go:embed templates/gitignore_editors.tpl
var gitignoreEditors string

//go:embed templates/gitignore_terraform.tpl
var gitignoreTerraform string

//go:embed templates/gitignore_helm.tpl
var gitignoreHelm string

//go:embed templates/helmignore.tpl
var helmignore string

func TerraformGitignore() string {
	return gitignoreTerraform + "\n" + gitignoreEditors
}

func HelmGitignore() string {
	return gitignoreEditors + "\n" + gitignoreHelm
}

func Helmignore() string {
	return helmignore
}
