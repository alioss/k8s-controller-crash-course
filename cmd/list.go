package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes deployments in the default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		// Configure logging
		level := parseLogLevel(logLevel)
		configureLogger(level)
		
		log.Debug().Str("kubeconfig", kubeconfig).Msg("Creating Kubernetes client")
		
		clientset, err := getKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}
		
		log.Debug().Msg("Fetching deployments from default namespace")
		
		deployments, err := clientset.AppsV1().Deployments("default").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to list deployments")
			os.Exit(1)
		}
		
		// Display header
		fmt.Printf("\n🎯 Kubernetes Deployments in 'default' namespace\n")
		fmt.Printf("═══════════════════════════════════════════════════\n")
		fmt.Printf("Found %d deployment(s):\n\n", len(deployments.Items))
		
		if len(deployments.Items) == 0 {
			fmt.Printf("❌ No deployments found in 'default' namespace\n\n")
			return
		}
		
		// Display each deployment with details
		for i, d := range deployments.Items {
			displayDeployment(d, i+1)
		}
		
		// Summary
		readyCount := 0
		totalReplicas := int32(0)
		readyReplicas := int32(0)
		
		for _, d := range deployments.Items {
			if d.Status.Replicas > 0 && d.Status.ReadyReplicas == d.Status.Replicas {
				readyCount++
			}
			totalReplicas += d.Status.Replicas
			readyReplicas += d.Status.ReadyReplicas
		}
		
		fmt.Printf("\n📊 Summary:\n")
		fmt.Printf("   Ready Deployments: %d/%d\n", readyCount, len(deployments.Items))
		fmt.Printf("   Total Pods: %d/%d running\n\n", readyReplicas, totalReplicas)
		
		log.Info().Int("total_deployments", len(deployments.Items)).Int("ready_deployments", readyCount).Msg("Successfully listed deployments")
	},
}

// displayDeployment shows detailed information about a single deployment
func displayDeployment(d appsv1.Deployment, index int) {
	// Get status indicator
	statusIcon := getStatusIcon(d)
	
	// Get age
	age := time.Since(d.CreationTimestamp.Time)
	ageStr := formatAge(age)
	
	// Get image (first container)
	image := "unknown"
	if len(d.Spec.Template.Spec.Containers) > 0 {
		image = d.Spec.Template.Spec.Containers[0].Image
	}
	
	// Display deployment info
	fmt.Printf("%s %d. %s\n", statusIcon, index, d.Name)
	fmt.Printf("   📦 Replicas: %d/%d ready", d.Status.ReadyReplicas, d.Status.Replicas)
	
	if d.Status.UpdatedReplicas != d.Status.Replicas {
		fmt.Printf(" (updating: %d)", d.Status.UpdatedReplicas)
	}
	
	fmt.Printf("\n   🏷️  Image: %s", image)
	fmt.Printf("\n   📅 Age: %s", ageStr)
	
	// Show conditions if any issues
	if len(d.Status.Conditions) > 0 {
		for _, condition := range d.Status.Conditions {
			if condition.Status != "True" && condition.Type == "Available" {
				fmt.Printf("\n   ⚠️  Issue: %s", condition.Message)
			}
		}
	}
	
	fmt.Printf("\n\n")
}

// getStatusIcon returns appropriate emoji based on deployment status
func getStatusIcon(d appsv1.Deployment) string {
	if d.Status.Replicas == 0 {
		return "⏸️"  // stopped
	}
	if d.Status.ReadyReplicas == d.Status.Replicas {
		return "✅"  // healthy
	}
	if d.Status.ReadyReplicas == 0 {
		return "❌"  // failed
	}
	return "⚠️"   // partial/updating
}

// formatAge formats duration into human readable string
func formatAge(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.0fh", d.Hours())
	}
	return fmt.Sprintf("%.0fd", d.Hours()/24)
}

func getKubeClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
}
