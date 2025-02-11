package kubernetes

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DeployImage creates (or replaces) the regcred secret in the target namespace
// and then deploys an App custom resource.
func (k *KubernetesUtil) DeployImage(repoName, image string) (*unstructured.Unstructured, error) {
	// --- Create or Replace the "regcred" Secret ---

	// Define the GroupVersionResource for Secrets.
	secretGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}

	secretName := "regcred"
	secretNamespace := k.config.DeploymentNamespace

	// Extract the registry host from the image string.
	// Expected image format: "registry-host/namespace/imagename:tag"
	parts := strings.Split(image, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("image name %q does not contain a registry hostname", image)
	}
	registryHost := parts[0]

	// Retrieve registry credentials from configuration or use defaults.
	username := k.config.ContainerRegistryUser
	password := k.config.ContainerRegistryUser

	// Create the docker config JSON for authenticating to the registry.
	// The structure follows:
	// {
	//   "auths": {
	//     "registryHost": {
	//       "username": "username",
	//       "password": "password",
	//       "auth": "<base64(username:password)>"
	//     }
	//   }
	// }
	authEncoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	dockerConfig := map[string]interface{}{
		"auths": map[string]interface{}{
			registryHost: map[string]string{
				"username": username,
				"password": password,
				"auth":     authEncoded,
			},
		},
	}
	dockerConfigBytes, err := json.Marshal(dockerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal docker config JSON: %v", err)
	}
	// The data value must be base64-encoded.
	dockerConfigBase64 := base64.StdEncoding.EncodeToString(dockerConfigBytes)

	// Create the secret object.
	secretObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      secretName,
				"namespace": secretNamespace,
			},
			"data": map[string]interface{}{
				".dockerconfigjson": dockerConfigBase64,
			},
			"type": "kubernetes.io/dockerconfigjson",
		},
	}

	// Create a context for API calls.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt to create the secret. If it already exists, update it.
	_, err = k.client.Resource(secretGVR).Namespace(secretNamespace).Create(ctx, secretObj, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			_, err = k.client.Resource(secretGVR).Namespace(secretNamespace).Update(ctx, secretObj, metav1.UpdateOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to update regcred secret: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to create regcred secret: %v", err)
		}
	}

	// --- Create the App Custom Resource ---

	// Define the GroupVersionResource for your custom resource.
	appGVR := schema.GroupVersionResource{
		Group:    "app.unbind.cloud",
		Version:  "v1",
		Resource: "apps", // plural name of your custom resource
	}

	// Create an unstructured object representing your custom resource.
	appCR := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "app.unbind.cloud/v1",
			"kind":       "App",
			"metadata": map[string]interface{}{
				"name":      repoName,
				"namespace": secretNamespace, // same namespace as the secret
			},
			"spec": map[string]interface{}{
				"image":           image,
				"domain":          fmt.Sprintf("%s.unbind.app", strings.ReplaceAll(repoName, "_", "-")),
				"imagePullSecret": "regcred",
			},
		},
	}

	// Create the custom resource in the target namespace.
	createdCR, err := k.client.Resource(appGVR).Namespace(secretNamespace).Create(ctx, appCR, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create custom resource: %v", err)
	}

	return createdCR, nil
}
