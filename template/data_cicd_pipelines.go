package template

type TerraformInfraTemplate struct {
	Project       string
	Repo          string
	Modules       []TerraformInfraModuleTemplate
	DefaultBranch string
}

type TerraformInfraModuleTemplate struct {
	Name         string
	Dependencies []string
}

type TerraformModuleTemplate struct {
	Project       string
	Repo          string
	Modules       []string
	DefaultBranch string
}
