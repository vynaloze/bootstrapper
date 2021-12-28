package template

type GitRepoType string

const (
	TerraformInfra  GitRepoType = "tf_infra"
	TerraformModule GitRepoType = "tf_module"
	Miscellaneous   GitRepoType = "misc"
)
