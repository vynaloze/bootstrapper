package template

type TerraformInfraTemplate struct {
	Project         string
	Repo            string
	Module          string
	Workflow        string
	DefaultBranch   string
	OnDefaultBranch bool
}
