package blueprint

import "bootstrapper/actor/terraform"

type TerraformOpts struct {
	terraform.Opts
	TerraformInfraCoreDir *string
}

var defaultTerraformOpts = TerraformOpts{
	TerraformInfraCoreDir: ptr("core"),
}

func (o *TerraformOpts) GetTerraformInfraCoreDir() string {
	if o.TerraformInfraCoreDir == nil {
		return *defaultTerraformOpts.TerraformInfraCoreDir
	}
	return *o.TerraformInfraCoreDir
}
