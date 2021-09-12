package template

import (
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

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
