package template

import (
	"bytes"
	_ "embed"
	"text/template"
)

func parse(tpl string, data interface{}) ([]byte, error) {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//TODO dynamic with provider
//go:embed templates/tf-infra-shared/core/repos_github.tf
var TfInfraSharedCoreReposTf string

//go:embed templates/tf-infra-shared/core/variables.tf
var TfInfraSharedCoreVariablesTf string

//TODO dynamic with provider
//go:embed templates/tf-infra-shared/core/versions_github.tf
var TfInfraSharedCoreVersionsTf string
