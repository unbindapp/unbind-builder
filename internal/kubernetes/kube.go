package kubernetes

import (
	"github.com/unbindapp/unbind-builder/config"
	"github.com/unbindapp/unbind-builder/internal/log"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
		if cfg.KubeConfig != "" {
			// Use the configured kubeconfig file instead
			kubeConfig, err = clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)
			if err != nil {
				log.Fatalf("Error building kubeconfig from %s: %v", cfg.KubeConfig, err)
			}
			log.Infof("Using kubeconfig from: %s", cfg.KubeConfig)
		} else {
			log.Fatalf("Error getting in-cluster config: %v", err)
		}
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
