package template

type K8sAppsTemplate struct {
	GitProvider      string
	GitProject       string
	GitRepo          string
	GitDefaultBranch string
	Environment      string
}

type K8sInfraArgoCdTemplate struct {
	GitProvider      string
	GitProject       string
	GitDefaultBranch string
	GitInfraRepo     string
	GitAppsRepo      string
	Environment      string
	Domain           string
}

type K8sInfraExternalDnsTemplate struct {
	InternalTxtOwnerId        string
	InternalDomain            string
	InternalDnsProvider       string
	InternalDnsProviderConfig map[string]string

	ServiceAccountAnnotations map[string]string
}

type K8sInfraIngressNginxTemplate struct {
	InternalLoadBalancerConfig map[string]string
}
