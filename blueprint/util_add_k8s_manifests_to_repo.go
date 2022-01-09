package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"log"
	"time"
)

type AddK8sManifestsToRepoOpts struct {
	TargetRepoOpts    git.Opts
	Templates         []Template
	TemplatesBasePath string
}

func AddK8sManifestsToRepo(opts *AddK8sManifestsToRepoOpts) error {
	log.Printf("adding k8s manifests to %s repo", opts.TargetRepoOpts.Repo)

	localActor := git.NewLocal(&opts.TargetRepoOpts)
	remoteActor, err := git.NewRemote(&opts.TargetRepoOpts)
	if err != nil {
		return fmt.Errorf("cannot initialize remote Git actor: %w", err)
	}

	log.Printf("preparing templates")

	gitFiles, err := templatesToGitFiles(opts.TemplatesBasePath, opts.Templates)
	if err != nil {
		return fmt.Errorf("error preparing templates: %w", err)
	}

	log.Printf("pushing changes to remote repository")

	message := "chore: add k8s manifests"
	branch := fmt.Sprintf("%s/%d", opts.TargetRepoOpts.GetAuthorName(), time.Now().UnixMilli())
	err = commitAndPush(localActor, branch, message, gitFiles)
	if err != nil {
		return err
	}
	err = remoteActor.RequestReview(&branch, &message)
	if err != nil {
		return fmt.Errorf("error creating PR: %w", err)
	}

	return nil
}

func K8sAppsPreset(k8sInfraGitOpts git.Opts, environment string) []Template {
	return []Template{
		{
			SourceFile: "apps/external-dns.yaml",
			Data: template.K8sAppsTemplate{
				GitProvider:      k8sInfraGitOpts.Provider,
				GitProject:       k8sInfraGitOpts.Project,
				GitRepo:          k8sInfraGitOpts.Repo,
				GitDefaultBranch: k8sInfraGitOpts.GetDefaultBranch(),
				Environment:      environment,
			},
			TargetFile: "apps/external-dns.yaml",
		},
		{
			SourceFile: "apps/ingress-nginx.yaml",
			Data: template.K8sAppsTemplate{
				GitProvider:      k8sInfraGitOpts.Provider,
				GitProject:       k8sInfraGitOpts.Project,
				GitRepo:          k8sInfraGitOpts.Repo,
				GitDefaultBranch: k8sInfraGitOpts.GetDefaultBranch(),
			},
			TargetFile: "apps/ingress-nginx.yaml",
		},
		{Source: template.HelmGitignore(), TargetFile: ".gitignore"},
	}
}

type K8sAwsInfraOpts struct {
	Environment        string
	Domain             string
	Route53ZoneId      string
	Region             string
	ExternalDnsRoleArn string
}

func K8sInfraAwsPreset(k8sInfraGitOpts git.Opts, k8sAppsGitOpts git.Opts, awsOpts K8sAwsInfraOpts) []Template {
	externalDnsData := template.K8sInfraExternalDnsTemplate{
		InternalTxtOwnerId:  awsOpts.Route53ZoneId,
		InternalDomain:      awsOpts.Domain,
		InternalDnsProvider: "aws",
		InternalDnsProviderConfig: map[string]string{
			"region":   awsOpts.Region,
			"zoneType": "private",
		},
		ServiceAccountAnnotations: map[string]string{
			"eks.amazonaws.com/role-arn": awsOpts.ExternalDnsRoleArn,
		},
	}
	return []Template{
		{Source: template.Helmignore(), TargetFile: "argocd/.helmignore"},
		{SourceFile: "argocd/Chart.lock", TargetFile: "argocd/Chart.lock"},
		{SourceFile: "argocd/Chart.yaml", TargetFile: "argocd/Chart.yaml"},
		{SourceFile: "argocd/values.yaml", TargetFile: "argocd/values.yaml"},
		{
			SourceFile: "argocd/values-env.yaml",
			Data: template.K8sInfraArgoCdTemplate{
				GitProvider:      k8sInfraGitOpts.Provider,
				GitProject:       k8sInfraGitOpts.Project,
				GitDefaultBranch: k8sInfraGitOpts.GetDefaultBranch(),
				GitInfraRepo:     k8sInfraGitOpts.Repo,
				GitAppsRepo:      k8sAppsGitOpts.Repo,
				Environment:      awsOpts.Environment,
				Domain:           awsOpts.Domain,
			},
			TargetFile: fmt.Sprintf("argocd/values-%s.yaml", awsOpts.Environment),
		},
		{Source: template.Helmignore(), TargetFile: "external-dns/.helmignore"},
		{SourceFile: "external-dns/Chart.lock", TargetFile: "external-dns/Chart.lock"},
		{SourceFile: "external-dns/Chart.yaml", TargetFile: "external-dns/Chart.yaml"},
		{
			SourceFile: "external-dns/values.yaml",
			Data:       externalDnsData,
			TargetFile: "external-dns/values.yaml",
		},
		{
			SourceFile: "external-dns/values-env.yaml",
			Data:       externalDnsData,
			TargetFile: fmt.Sprintf("external-dns/values-%s.yaml", awsOpts.Environment),
		},
		{Source: template.Helmignore(), TargetFile: "ingress-nginx/.helmignore"},
		{SourceFile: "ingress-nginx/Chart.lock", TargetFile: "ingress-nginx/Chart.lock"},
		{SourceFile: "ingress-nginx/Chart.yaml", TargetFile: "ingress-nginx/Chart.yaml"},
		{
			SourceFile: "ingress-nginx/values.yaml",
			Data: template.K8sInfraIngressNginxTemplate{
				InternalLoadBalancerConfig: map[string]string{
					"service.beta.kubernetes.io/aws-load-balancer-internal": "0.0.0.0/0",
					"service.beta.kubernetes.io/aws-load-balancer-type":     "nlb",
				},
			},
			TargetFile: "ingress-nginx/values.yaml",
		},
		{Source: template.HelmGitignore(), TargetFile: ".gitignore"},
	}
}
