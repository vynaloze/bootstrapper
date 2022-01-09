package helm

import (
	"context"
	"fmt"
	"github.com/vynaloze/go-helm-client"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/repo"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Opts struct {
	KubeConfigPath string
}

type Actor interface {
	Install(opts InstallOpts) error
}

type InstallOpts struct {
	Name        string
	Path        string
	ValuesFiles []string
	// secret name => (key => value)
	ExtraSecrets map[string]map[string][]byte

	Namespace *string
}

type actor struct {
	helmClient helmclient.Client
	k8sClient  *kubernetes.Clientset
}

func New(opts Opts) (Actor, error) {
	if opts.KubeConfigPath == "" {
		home := os.Getenv("HOME")
		opts.KubeConfigPath = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", opts.KubeConfigPath)
	if err != nil {
		log.Panicln("failed to create k8s config")
	}

	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create k8s client")
	}

	kubeConfig, err := ioutil.ReadFile(opts.KubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig file at %s: %w", opts.KubeConfigPath, err)
	}

	opt := &helmclient.KubeConfClientOptions{
		Options: &helmclient.Options{
			DebugLog: func(format string, v ...interface{}) {
				msg := fmt.Sprintf(format, v...)
				if strings.Contains(msg, "is not ready") {
					return
				}
				log.Printf("[k8s] %s", msg)
			},
		},
		KubeConfig: kubeConfig,
	}

	hc, err := helmclient.NewClientFromKubeConf(opt)
	if err != nil {
		return nil, fmt.Errorf("cannot create Helm client")
	}

	return &actor{hc, kc}, nil
}

func (a *actor) Install(opts InstallOpts) error {
	if opts.Namespace == nil {
		n := "default"
		opts.Namespace = &n
	}

	log.Printf("installing Helm repos")
	if err := a.addChartDeps(opts); err != nil {
		return fmt.Errorf("cannot add Helm repos: %w", err)
	}

	log.Printf("creating secrets")
	if err := a.createSecrets(opts); err != nil {
		return fmt.Errorf("cannot create secrets: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current direcrory: %w", err)
	}
	defer os.Chdir(cwd)
	err = os.Chdir(opts.Path)
	if err != nil {
		return fmt.Errorf("cannot change to chart directory (%s): %w", opts.Path, err)
	}

	valuesYaml := ""
	for _, vf := range opts.ValuesFiles {
		c, err := ioutil.ReadFile(filepath.Join(opts.Path, vf))
		if err != nil {
			return fmt.Errorf("falied to read values file %s: %w", vf, err)
		}
		valuesYaml = valuesYaml + "\n" + string(c)
	}

	chartSpec := helmclient.ChartSpec{
		ReleaseName:      opts.Name,
		ChartName:        ".",
		Namespace:        *opts.Namespace,
		ValuesYaml:       valuesYaml,
		DependencyUpdate: true,
		Wait:             true,
		Timeout:          5 * time.Minute,
	}

	log.Printf("update chart dependencies")
	if err := a.helmClient.UpdateDependencies(&chartSpec); err != nil {
		return fmt.Errorf("cannot download Helm chart dependencies: %w", err)
	}
	log.Printf("installing release")
	if _, err := a.helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec); err != nil {
		return fmt.Errorf("cannot install Helm chart: %w", err)
	}
	return nil
}

func (a *actor) addChartDeps(opts InstallOpts) error {
	chartDepsContent, err := ioutil.ReadFile(filepath.Join(opts.Path, "Chart.yaml"))
	if err != nil {
		return fmt.Errorf("cannot read chart dependencies: %w", err)
	}
	var chartDepsYaml chartDeps
	if err := yaml.Unmarshal(chartDepsContent, &chartDepsYaml); err != nil {
		return fmt.Errorf("cannot read chart dependencies: %w", err)
	}
	for _, dep := range chartDepsYaml.Dependencies {
		chartRepo := repo.Entry{
			Name: dep.Name,
			URL:  dep.Repository,
		}
		if err := a.helmClient.AddOrUpdateChartRepo(chartRepo); err != nil {
			return fmt.Errorf("cannot add Helm repo %s (%s): %w", dep.Name, dep.Repository, err)
		}
	}
	return nil
}

type chartDeps struct {
	Dependencies []chartDep `yaml:"dependencies"`
}
type chartDep struct {
	Name       string `yaml:"name"`
	Repository string `yaml:"repository"`
}

func (a *actor) createSecrets(opts InstallOpts) error {
	if len(opts.ExtraSecrets) == 0 {
		return nil
	}

	secretsClient := a.k8sClient.CoreV1().Secrets(*opts.Namespace)

	for secretName, secretData := range opts.ExtraSecrets {
		secret := v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: *opts.Namespace,
			},
			Type: v1.SecretTypeOpaque,
			Data: secretData,
		}

		if _, err := secretsClient.Create(context.Background(), &secret, metav1.CreateOptions{}); err != nil {
			if k8serrors.IsAlreadyExists(err) {
				if _, err = secretsClient.Update(context.Background(), &secret, metav1.UpdateOptions{}); err != nil {
					return fmt.Errorf("cannot update secret %s : %w", secretName, err)
				}
			} else {
				return fmt.Errorf("cannot create secret %s : %w", secretName, err)
			}
		}

	}
	return nil
}
