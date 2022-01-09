package main

import (
	"bootstrapper/actor/git"
	"bootstrapper/actor/helm"
	"bootstrapper/blueprint"
	"bootstrapper/datasource"
	"bootstrapper/template"
	"fmt"
	"log"
)

func main() {
	secrets, err := datasource.NewYamlFile("secrets.yaml")
	if err != nil {
		panic(err)
	}
	ghToken, ok := secrets.Get("actor.git.github.token")
	if !ok {
		panic(err)
	}
	sshKey, ok := secrets.Get("actor.git.github.ssh-private-key")
	if !ok {
		panic(err)
	}

	gitProvider := "github.com"
	gitProject := "bootstrapper-demo-org"
	gitUser := "bootstrapper-demo"
	gitPass := ghToken

	env := "stg"
	region := "eu-central-1"
	domain := "stg.euc1.aws.b-demo.org"
	r53ZoneId := "XXX"
	extDnsRoleArn := "arn:aws:iam::XXX:role/stg-external-dns"

	sharedInfraGitOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "tf-infra-shared",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}
	k8sInfraRepoOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "k8s-infra",
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}
	k8sAppsRepoOpts := git.Opts{
		Provider: gitProvider, Project: gitProject, Repo: "k8s-apps-" + env,
		RemoteAuthUser: gitUser, RemoteAuthPass: gitPass,
	}

	helmActor, err := helm.New(helm.Opts{})
	if err != nil {
		panic(err)
	}

	log.Printf("start k8s phase")

	log.Printf("create repos")
	createGitRepoOpts := blueprint.CreateGitReposOpts{
		SharedInfraRepoOpts: sharedInfraGitOpts,
		NewReposSpecs: []blueprint.CreateGitReposNewRepoSpec{
			{NewRepoOpts: k8sInfraRepoOpts, NewRepoType: template.Miscellaneous},
			{NewRepoOpts: k8sAppsRepoOpts, NewRepoType: template.Miscellaneous},
		},
	}
	if err := blueprint.CreateGitRepos(createGitRepoOpts); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Press Enter to proceed")
	fmt.Scanln()

	log.Printf("setup repos")
	k8sAppsManifestsOpts := &blueprint.AddK8sManifestsToRepoOpts{
		TargetRepoOpts:    k8sAppsRepoOpts,
		Templates:         blueprint.K8sAppsPreset(k8sInfraRepoOpts, env),
		TemplatesBasePath: "k8s-apps",
	}
	k8sInfraManifestsOpts := &blueprint.AddK8sManifestsToRepoOpts{
		TargetRepoOpts: k8sInfraRepoOpts,
		Templates: blueprint.K8sInfraAwsPreset(k8sInfraRepoOpts, k8sAppsRepoOpts, blueprint.K8sAwsInfraOpts{
			Environment:        env,
			Domain:             domain,
			Route53ZoneId:      r53ZoneId,
			Region:             region,
			ExternalDnsRoleArn: extDnsRoleArn,
		}),
		TemplatesBasePath: "k8s-infra",
	}
	if err := blueprint.AddK8sManifestsToRepo(k8sAppsManifestsOpts); err != nil {
		log.Fatalln(err)
	}
	if err := blueprint.AddK8sManifestsToRepo(k8sInfraManifestsOpts); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Press Enter to proceed")
	fmt.Scanln()

	log.Printf("install ArgoCD")
	k8sInfraLocalRepoPath, err := git.NewLocal(&k8sInfraRepoOpts).Clone()
	if err != nil {
		log.Fatalln(err)
	}
	installOpts := helm.InstallOpts{
		Name:        "argocd",
		Path:        k8sInfraLocalRepoPath + "/argocd",
		ValuesFiles: []string{fmt.Sprintf("values-%s.yaml", env)},
		ExtraSecrets: map[string]map[string][]byte{
			"argocd-git-ssh-readonly-k8s-infra": {
				"private-key": []byte(sshKey),
			},
			"argocd-git-ssh-readonly-k8s-apps-" + env: {
				"private-key": []byte(sshKey),
			},
		},
	}
	err = helmActor.Install(installOpts)
	if err != nil {
		log.Fatalln(err)
	}
}
