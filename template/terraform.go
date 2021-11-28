package template

import (
	_ "embed"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"regexp"
	"sort"
)

//go:embed templates/gitignore_editors.tpl
var gitignoreEditors string

//go:embed templates/gitignore_terraform.tpl
var gitignoreTerraform string

func TerraformGitignore() string {
	return gitignoreTerraform + "\n" + gitignoreEditors
}

type TerraformVariable struct {
	Key   string
	Value cty.Value
}

func TerraformModuleFromRegistry(name string, source string, version string, vars []*TerraformVariable) (string, error) {
	// vars has to be list to preserve variable ordering and allow specifying whitespace in between
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	moduleBlock := rootBody.AppendNewBlock("module", []string{name})
	moduleBody := moduleBlock.Body()
	moduleBody.SetAttributeValue("source", cty.StringVal(source))
	moduleBody.SetAttributeValue("version", cty.StringVal(version))
	moduleBody.AppendNewline()
	for _, v := range vars {
		if v == nil {
			moduleBody.AppendNewline()
		} else {
			//t, err := gocty.ImpliedType(v.Value)
			//if err != nil {
			//	return "", fmt.Errorf("gocty.ImpliedType(%v): %w", v.Value, err)
			//}
			//vv, err := gocty.ToCtyValue(v.Value, t)
			//if err != nil {
			//	return "", fmt.Errorf("gocty.ToCtyValue(%v, %v): %w", v.Value, t, err)
			//}
			moduleBody.SetAttributeValue(v.Key, v.Value)
		}
	}
	rootBody.AppendNewline()
	return string(f.Bytes()), nil
}

type TerraformProvider struct {
	Source  string
	Version string
}

func TerraformVersionsFile(providers map[string]TerraformProvider) (string, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	terraformBlock := rootBody.AppendNewBlock("terraform", nil)
	terraformBody := terraformBlock.Body()
	providersBlock := terraformBody.AppendNewBlock("required_providers", nil)
	providersBody := providersBlock.Body()
	for _, name := range sortedKeys(providers) {
		// There should be just:
		// providersBody.SetAttributeValue(name, cty.ObjectVal(map[string]cty.Value{...}))
		// but hclwrite is not able to produce a map where every entry is on a new line,
		// instead it puts entire map in a single line. This is a workaround for pretty-formatting:
		// produce block now and later just add "=" between block name and "{" (using blocksAsMap() function).
		providerBlock := providersBody.AppendNewBlock(name, nil)
		providerBody := providerBlock.Body()
		providerBody.SetAttributeValue("source", cty.StringVal(providers[name].Source))
		providerBody.SetAttributeValue("version", cty.StringVal(providers[name].Version))
	}
	return blocksAsMap(string(f.Bytes()), providers)
}

func mapAsBlock(sourceMap map[string]cty.Value) {
	// try with go templates, maybe?
}

func blocksAsMap(body string, providers map[string]TerraformProvider) (string, error) {
	for name := range providers {
		re, err := regexp.Compile(name + " {")
		if err != nil {
			return "", err
		}
		body = re.ReplaceAllString(body, name+" = {")
	}
	return body, nil
}

func sortedKeys(m map[string]TerraformProvider) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
