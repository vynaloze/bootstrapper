package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Opts struct {
	TerraformCloudOrg       string
	TerraformCloudToken     string
	TerraformCloudWorkspace string

	ProviderSecrets map[string]map[string]string
}

type Actor struct {
	execPath string
	opts     *Opts
}

func New(opts *Opts) (*Actor, error) {
	var execPath string
	execPath, err := tfinstall.Find(context.Background(), tfinstall.LookPath())
	if err != nil {
		log.Println("Terraform binary not found on PATH. Installing a fresh one")
		installer := &releases.ExactVersion{
			Product: product.Terraform,
			Version: version.Must(version.NewVersion("1.0.11")),
		}

		execPath, err = installer.Install(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error installing Terraform: %w", err)
		}
	}
	return &Actor{execPath, opts}, nil
}

func (a *Actor) Apply(dir string) error {
	tf, err := tfexec.NewTerraform(dir, a.execPath)
	if err != nil {
		return fmt.Errorf("error running Terraform: %w", err)
	}
	tf.SetLogger(log.Default())

	err = a.writeCliConfigFile(dir)
	if err != nil {
		return fmt.Errorf("error running Terraform: %w", err)
	}
	err = a.writeProvidersTf(dir)
	if err != nil {
		return fmt.Errorf("error running Terraform: %w", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return fmt.Errorf("error running Init: %w", err)
	}
	wsId, err := a.getWorkspaceId()
	if err != nil {
		return fmt.Errorf("error fetching workspace ID: %w", err)
	}
	err = tf.Import(context.Background(), fmt.Sprintf("tfe_workspace.this[\"%s\"]", a.opts.TerraformCloudWorkspace), wsId)
	if err != nil && !strings.Contains(err.Error(), "Resource already managed by Terraform") {
		return fmt.Errorf("error running Import: %w", err)
	}

	err = tf.Apply(context.Background())
	if err != nil {
		return fmt.Errorf("error running Apply: %w", err)
	}

	return nil
}

func (a *Actor) getWorkspaceId() (string, error) {
	url := fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%s/workspaces/%s",
		a.opts.TerraformCloudOrg, a.opts.TerraformCloudWorkspace)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error preparing HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization", "Bearer "+a.opts.TerraformCloudToken)
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading reponse: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got %d response: %s", resp.StatusCode, body)
	}
	var jsonResp showWorkspaceResponse
	if err = json.Unmarshal(body, &jsonResp); err != nil {
		return "", fmt.Errorf("error decoding reponse: %w", err)
	}
	return jsonResp.Data.Id, nil
}

type showWorkspaceResponse struct {
	Data showWorkspaceResponseData `json:"data"`
}
type showWorkspaceResponseData struct {
	Id string `json:"id"`
}
