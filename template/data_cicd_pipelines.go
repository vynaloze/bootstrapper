package template

type CICDTerraformInfraTemplate struct {
	Project       string
	Repo          string
	Modules       []CICDTerraformInfraModuleTemplate
	DefaultBranch string
}

type CICDTerraformInfraModuleTemplate struct {
	Name         string
	Dependencies []string
}

type CICDTerraformModuleTemplate struct {
	Project       string
	Repo          string
	Modules       []string
	DefaultBranch string
}
