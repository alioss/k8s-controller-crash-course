package client

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"atlasctl/pkg/models"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// K8sClient wraps the Kubernetes client
type K8sClient struct {
	clientset *kubernetes.Clientset
}

// NewK8sClient creates a new Kubernetes client with default kubeconfig
func NewK8sClient() (*K8sClient, error) {
	return NewK8sClientWithConfig("")
}

// NewK8sClientWithConfig creates a new Kubernetes client with specified kubeconfig
func NewK8sClientWithConfig(kubeconfigPath string) (*K8sClient, error) {
	config, err := getConfigWithPath(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &K8sClient{clientset: clientset}, nil
}

// getConfigWithPath returns kubernetes config from specified kubeconfig path or default
func getConfigWithPath(kubeconfigPath string) (*rest.Config, error) {
	// Try in-cluster config first if no specific kubeconfig provided
	if kubeconfigPath == "" {
		if config, err := rest.InClusterConfig(); err == nil {
			return config, nil
		}
	}

	// Determine kubeconfig path
	var kubeconfig string
	if kubeconfigPath != "" {
		// Use provided path
		kubeconfig = kubeconfigPath
	} else {
		// Fall back to default kubeconfig
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig %s: %w", kubeconfig, err)
	}

	return config, nil
}

// GetAtlasApps retrieves Atlas applications from specified namespaces
func (k *K8sClient) GetAtlasApps(ctx context.Context, namespaces []string) (models.AtlasAppList, error) {
	var apps models.AtlasAppList

	for _, ns := range namespaces {
		deployments, err := k.clientset.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{
			LabelSelector: "app=atlas", // Filter for atlas apps
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", ns, err)
		}

		for _, deployment := range deployments.Items {
			app := k.convertDeploymentToAtlasApp(deployment)
			apps = append(apps, app)
		}
	}

	return apps, nil
}

// convertDeploymentToAtlasApp converts a Kubernetes Deployment to AtlasApp
func (k *K8sClient) convertDeploymentToAtlasApp(deployment appsv1.Deployment) models.AtlasApp {
	app := models.AtlasApp{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
	}

	// Extract version from image
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		container := deployment.Spec.Template.Spec.Containers[0]
		app.Version = extractVersionFromImage(container.Image)

		// Extract migration ID from environment variables
		app.MigrationID = extractMigrationID(container.Env)
	}

	// Set status based on replica status
	ready := deployment.Status.ReadyReplicas
	total := *deployment.Spec.Replicas
	app.Status = string(models.GetStatus(ready, total))
	app.Replicas = fmt.Sprintf("%d/%d", ready, total)

	// Set last update time
	if len(deployment.Status.Conditions) > 0 {
		app.LastUpdate = deployment.Status.Conditions[len(deployment.Status.Conditions)-1].LastUpdateTime.Time
	}
	if app.LastUpdate.IsZero() {
		app.LastUpdate = deployment.CreationTimestamp.Time
	}

	// Calculate age
	app.Age = formatAge(time.Since(deployment.CreationTimestamp.Time))

	return app
}

// extractVersionFromImage extracts version from container image
func extractVersionFromImage(image string) string {
	// Extract version from image like "nginx:1.21.0" -> "1.21.0"
	parts := strings.Split(image, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return "unknown"
}

// extractMigrationID extracts migration ID from environment variables
func extractMigrationID(envVars []corev1.EnvVar) string {
	for _, env := range envVars {
		if env.Name == "MIGRATION_ID" {
			return env.Value
		}
	}
	return "unknown"
}

// formatAge formats duration to human readable string
func formatAge(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	} else {
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
