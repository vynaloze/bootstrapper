package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"log"
)

type Actor struct {
	execPath string
}

func New() (*Actor, error) {
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
	return &Actor{execPath}, nil
}

func (a *Actor) Apply(dir string) error {
	tf, err := tfexec.NewTerraform(dir, a.execPath)
	if err != nil {
		return fmt.Errorf("error running Terraform: %w", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return fmt.Errorf("error running Init: %w", err)
	}

	err = tf.Apply(context.Background())
	if err != nil {
		return fmt.Errorf("error running Apply: %w", err)
	}

	return nil
}
