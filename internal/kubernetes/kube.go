package kubernetes

import (
	"log"

	"github.com/unbindapp/unbind-builder/config"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type KubernetesUtil struct {
	config    *config.Config
	namespace string
	client    *dynamic.DynamicClient
}

func NewKubernetesUtil(cfg *config.Config) *KubernetesUtil {
	// Get config
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting in-cluster config: %v", err)
	}

	clientset, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	return &KubernetesUtil{
		config:    cfg,
		client:    clientset,
		namespace: cfg.DeploymentNamespace,
	}
}
