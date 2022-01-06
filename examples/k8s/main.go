package main

import (
	"bootstrapper/actor/helm"
	"bootstrapper/datasource"
	"log"
)

func main() {
	secrets, err := datasource.NewYamlFile("secrets.yaml")
	if err != nil {
		panic(err)
	}
	sshKey, ok := secrets.Get("actor.git.github.ssh-private-key")
	if !ok {
		panic(err)
	}

	helmActor, err := helm.New(helm.Opts{})
	if err != nil {
		panic(err)
	}

	log.Printf("start k8s phase")

	log.Printf("setup repos")

	log.Printf("install ArgoCD")
	installOpts := helm.InstallOpts{
		Name:        "argocd",
		Path:        "/mnt/c/workspace/bootstrapper/template/templates/k8s-infra/argocd", //TODO
		ValuesFiles: []string{"values-dev.yaml"},
		ExtraSecrets: map[string]map[string][]byte{
			"argocd-git-ssh-readonly": {
				"private-key": []byte(sshKey),
			},
		},
	}
	err = helmActor.Install(installOpts)
	if err != nil {
		panic(err)
	}

}
