package blueprint

type TerraformOpts struct {
	ProviderSecrets map[string]string

	SharedInfraRepoName      *string
	TerraformInfraReposFile  *string
	TerraformModuleReposFile *string
}

var defaultTerraformOpts = TerraformOpts{
	SharedInfraRepoName:      ptr("tf-infra-shared"),
	TerraformInfraReposFile:  ptr("core/repos_tf_infra.tf"),
	TerraformModuleReposFile: ptr("core/repos_tf_module.tf"),
}

func (o *TerraformOpts) GetSharedInfraRepoName() string {
	if o.SharedInfraRepoName == nil {
		return *defaultTerraformOpts.SharedInfraRepoName
	}
	return *o.SharedInfraRepoName
}

func (o *TerraformOpts) GetTerraformInfraReposFile() string {
	if o.TerraformInfraReposFile == nil {
		return *defaultTerraformOpts.TerraformInfraReposFile
	}
	return *o.TerraformInfraReposFile
}

func (o *TerraformOpts) GetTerraformModuleReposFile() string {
	if o.TerraformModuleReposFile == nil {
		return *defaultTerraformOpts.TerraformModuleReposFile
	}
	return *o.TerraformModuleReposFile
}
