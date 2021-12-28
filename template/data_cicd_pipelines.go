package template

type TerraformInfraTemplate struct {
	Project       string
	Repo          string
	Module        string
	DefaultBranch string
}

type TerraformModuleTemplate struct {
	Project       string
	Repo          string
	Modules       []string
	DefaultBranch string
}
