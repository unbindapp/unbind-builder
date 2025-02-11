package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (k *KubernetesUtil) DeployImage(repoName, image string) (*unstructured.Unstructured, error) {
	// Define the GroupVersionResource for your CRD.
	appGVR := schema.GroupVersionResource{
		Group:    "app.unbind.cloud",
		Version:  "v1",
		Resource: "apps", // the plural name of your custom resource
	}

	// Create an unstructured object representing your custom resource.
	// Note that metadata.namespace is set to "unbind-user".
	appCR := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "app.unbind.cloud/v1",
			"kind":       "App",
			"metadata": map[string]interface{}{
				"name":      repoName,                     // name of the resource
				"namespace": k.config.DeploymentNamespace, // target namespace for the resource
			},
			"spec": map[string]interface{}{
				"image": image,
				// ! TODO - this should be a dynamic value
				"domain": fmt.Sprintf("%s.unbind.app", strings.ReplaceAll(repoName, "_", "-")),
			},
		},
	}

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create the custom resource in the target namespace
	createdCR, err := k.client.Resource(appGVR).Namespace(k.config.DeploymentNamespace).Create(ctx, appCR, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create custom resource: %v", err)
	}

	return createdCR, nil
}
