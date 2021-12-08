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

type TfInfraSharedCoreReposTfOpts struct {
	Repos         []string
	Strict        bool
	DefaultBranch string
}

//go:embed templates/tf-infra-shared/core/repos.tf.tpl
var tfInfraSharedCoreReposTf string

func TfInfraSharedCoreReposTf(data TfInfraSharedCoreReposTfOpts) (string, error) {
	parsed, err := parse(tfInfraSharedCoreReposTf, data)
	if err != nil {
		return "", err
	}
	return string(parsed), nil
}
