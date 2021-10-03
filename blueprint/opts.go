package blueprint

type TerraformOpts struct {
	SharedInfraRepoName     *string
	SharedInfraBootstrapDir *string
	SharedInfraCoreDir      *string
}

var defaultTerraformOpts = TerraformOpts{
	SharedInfraRepoName:     ptr("tf-infra-shared"),
	SharedInfraBootstrapDir: ptr("bootstrap"),
	SharedInfraCoreDir:      ptr("core"),
}

func (o *TerraformOpts) GetSharedInfraRepoName() string {
	if o.SharedInfraRepoName == nil {
		return *defaultTerraformOpts.SharedInfraRepoName
	}
	return *o.SharedInfraRepoName
}

func (o *TerraformOpts) GetSharedInfraBootstrapDir() string {
	if o.SharedInfraBootstrapDir == nil {
		return *defaultTerraformOpts.SharedInfraBootstrapDir
	}
	return *o.SharedInfraBootstrapDir
}

func (o *TerraformOpts) GetSharedInfraCoreDir() string {
	if o.SharedInfraCoreDir == nil {
		return *defaultTerraformOpts.SharedInfraCoreDir
	}
	return *o.SharedInfraCoreDir
}
