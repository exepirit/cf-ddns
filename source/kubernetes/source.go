package kubernetes

import (
	"context"
	"path/filepath"

	"github.com/exepirit/cf-ddns/domain"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func New() (*Source, error) {
	source, err := NewForServiceAccount()
	if err == nil {
		return source, nil
	}

	return NewForConfig()
}

func NewForConfig() (*Source, error) {
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	return &Source{
		client: clientset,
	}, err
}

func NewForServiceAccount() (*Source, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	return &Source{
		client: clientset,
	}, err
}

type Source struct {
	client *kubernetes.Clientset
}

func (s *Source) GetEndpoints() ([]*domain.Endpoint, error) {
	ctx := context.Background()

	domains, err := s.getDomainNames(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get domains")
	}

	targets, err := s.getTargets(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get targets")
	}

	endpoints := make([]*domain.Endpoint, len(domains))
	for i, d := range domains {
		endpoints[i] = domain.NewEndpoint(d, domain.RecordTypeA, targets...)
	}
	return endpoints, nil
}

func (s *Source) getTargets(ctx context.Context) (domain.Target, error) {
	nodesClient := s.client.CoreV1().Nodes()
	nodes, err := nodesClient.List(ctx, meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	target := make(domain.Target, 0)
	for _, node := range nodes.Items {
		addr := s.getNodeAddress(node)
		if addr != "" {
			target = append(target, addr)
		}
	}
	return target, nil
}

func (*Source) getNodeAddress(node core.Node) string {
	addresses := node.Status.Addresses

	for _, addr := range addresses {
		if addr.Type == core.NodeExternalIP {
			return addr.Address
		}
	}

	return ""
}

func (s *Source) getDomainNames(ctx context.Context) ([]string, error) {
	ingresses, err := s.getIngresses(ctx)
	if err != nil {
		return nil, err
	}

	domains := make([]string, 0)
	for _, ingress := range ingresses {
		for _, rule := range ingress.Spec.Rules {
			domains = append(domains, rule.Host)
		}
	}
	return domains, nil
}

func (s *Source) getIngresses(ctx context.Context) ([]networking.Ingress, error) {
	ingressesClient := s.client.NetworkingV1beta1().Ingresses("default")
	ingresses, err := ingressesClient.List(ctx, meta.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ingresses.Items, nil
}
