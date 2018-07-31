package prodruntime

import (
	"fmt"
	"net/url"

	restclient "k8s.io/client-go/rest"

	"github.com/bitnami/kube-prod-runtime/kubeprod/pkg/aks"
)

type Platform struct {
	Name        string
	Description string
	PreUpdate   func(config interface{}) (interface{}, error)
	PostUpdate  func(conf *restclient.Config) error
}

var Platforms = []Platform{
	{
		Name:        "minikube-0.25+k8s-1.9",
		Description: "Minikube 0.25 with Kubernetes 1.9",
	},
	{
		Name:        "minikube-0.25+k8s-1.8",
		Description: "Minikube 0.25 with Kubernetes 1.8",
	},
	{
		Name:        "aks+k8s-1.9",
		Description: "Azure Container Service (AKS) with Kubernetes 1.9",
		PreUpdate:   aks.PreUpdate,
	},
	{
		Name:        "aks+k8s-1.8",
		Description: "Azure Container Service (AKS) with Kubernetes 1.8",
		PreUpdate:   aks.PreUpdate,
	},
}

func FindPlatform(name string) *Platform {
	for i := range Platforms {
		p := &Platforms[i]
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (p *Platform) ManifestURL(base *url.URL) (*url.URL, error) {
	return base.Parse(fmt.Sprintf("platforms/%s.jsonnet", p.Name))
}

func (p *Platform) RunPreUpdate(kubeprodConf interface{}) (interface{}, error) {
	if p.PreUpdate == nil {
		return kubeprodConf, nil
	}
	return p.PreUpdate(kubeprodConf)
}

func (p *Platform) RunPostUpdate(conf *restclient.Config) error {
	if p.PostUpdate == nil {
		return nil
	}
	return p.PostUpdate(conf)
}