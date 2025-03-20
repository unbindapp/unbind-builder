package kubernetes

import (
	"context"
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DeployImage creates (or replaces) the service resource in the target namespace
// for deployment after a successful build job.
func (self *KubernetesUtil) DeployImage(ctx context.Context, repoName, image string) (*unstructured.Unstructured, error) {
	// Extract GitHub repository name from the Git URL
	// Typically in the format: https://github.com/org/repo.git
	gitRepoURL := self.config.GitRepoURL
	parts := strings.Split(gitRepoURL, "/")
	gitRepository := ""
	if len(parts) >= 2 {
		repoWithGit := parts[len(parts)-1]
		// Remove .git extension if present
		gitRepository = strings.TrimSuffix(repoWithGit, ".git")
		if len(parts) >= 3 {
			gitRepository = parts[len(parts)-2] + "/" + gitRepository
		}
	}

	// Generate a sanitized service name from the repo name
	serviceName := strings.ToLower(strings.ReplaceAll(repoName, "_", "-"))

	// Define the GroupVersionResource for the Service custom resource
	serviceGVR := schema.GroupVersionResource{
		Group:    "unbind.unbind.app",
		Version:  "v1",
		Resource: "services", // plural name of the custom resource
	}

	// Build service configuration
	serviceConfig := map[string]interface{}{
		// Git deployment settings
		"gitBranch":  self.config.GitRef,
		"autoDeploy": true,

		// Image from build
		"image": image,
	}

	if self.config.ServicePublic != nil {
		serviceConfig["public"] = *self.config.ServicePublic
	}

	if self.config.ServicePort != nil {
		serviceConfig["port"] = *self.config.ServicePort
	}

	if self.config.ServiceHost != nil {
		serviceConfig["host"] = *self.config.ServiceHost
	}

	if self.config.ServiceReplicas != nil {
		serviceConfig["replicas"] = *self.config.ServiceReplicas
	}

	// Create an unstructured object representing the Service custom resource
	serviceCR := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "unbind.unbind.app/v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name":      serviceName,
				"namespace": self.config.DeploymentNamespace,
			},
			"spec": map[string]interface{}{
				// Service identification
				"name":        serviceName,
				"displayName": serviceName,
				"description": fmt.Sprintf("Auto-deployed service for %s", repoName),

				// Service classification
				"type":      "git",
				"builder":   "railpack",
				"runtime":   self.config.ServiceRuntime,
				"framework": self.config.ServiceFramework,

				// Relations - may need to be configured appropriately for your environment
				"teamRef":       "default",
				"projectRef":    "default",
				"environmentId": "prod",

				// GitHub integration
				"githubInstallationId": self.config.GithubInstallationID,
				"gitRepository":        gitRepository,

				// Secret reference - creates a standard naming convention
				"kubernetesSecret": self.config.ServiceSecretName,

				// Configuration
				"config": serviceConfig,
			},
		},
	}

	// Create the custom resource in the target namespace
	createdCR, err := self.client.Resource(serviceGVR).Namespace(self.config.DeploymentNamespace).Create(ctx, serviceCR, metav1.CreateOptions{})
	if err != nil {
		// If the resource already exists, update it
		if apierrors.IsAlreadyExists(err) {
			// Retrieve the existing resource
			existingCR, getErr := self.client.Resource(serviceGVR).Namespace(self.config.DeploymentNamespace).Get(ctx, serviceCR.GetName(), metav1.GetOptions{})
			if getErr != nil {
				return nil, fmt.Errorf("failed to retrieve existing service: %v", getErr)
			}

			// Set the resourceVersion on the object to be updated
			serviceCR.SetResourceVersion(existingCR.GetResourceVersion())

			// Update the CR
			updatedCR, updateErr := self.client.Resource(serviceGVR).Namespace(self.config.DeploymentNamespace).Update(ctx, serviceCR, metav1.UpdateOptions{})
			if updateErr != nil {
				return nil, fmt.Errorf("failed to update service custom resource: %v", updateErr)
			}

			return updatedCR, nil
		}
		return nil, fmt.Errorf("failed to create service custom resource: %v", err)
	}

	return createdCR, nil
}
