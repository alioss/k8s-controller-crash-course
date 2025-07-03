package cmd

import (
	"context"
	"fmt"
	"time"

	"atlasctl/pkg/client"
	"atlasctl/pkg/formatter"

	"github.com/spf13/cobra"
)

var (
	namespaces   []string
	allNs        bool
	kubeconfigPath string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Atlas applications across environments",
	Long: `List Atlas applications running in different environments (dev, stage, prod).
Shows version, migration ID, status, and other deployment information.

Examples:
  # List Atlas apps in all default environments
  atlasctl list

  # List Atlas apps in specific namespaces
  atlasctl list --namespace dev,stage

  # List Atlas apps in all namespaces
  atlasctl list --all-namespaces`,
	
	RunE: func(cmd *cobra.Command, args []string) error {
		return runListCommand()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&namespaces, "namespace", "n", 
		[]string{"dev", "stage", "prod"}, 
		"Namespaces to search for Atlas apps (comma-separated)")
	
	listCmd.Flags().BoolVarP(&allNs, "all-namespaces", "A", false, 
		"Search in all namespaces")
	
	listCmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", "", 
		"Path to the kubeconfig file")
}

func runListCommand() error {
	fmt.Println("🔍 Searching for Atlas applications...")

	// Show kubeconfig info
	if kubeconfigPath != "" {
		fmt.Printf("🔗 Using kubeconfig: %s\n", kubeconfigPath)
	} else {
		fmt.Println("🔗 Using default kubeconfig")
	}

	// Create Kubernetes client with kubeconfig
	k8sClient, err := client.NewK8sClientWithConfig(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("✅ Successfully connected to Kubernetes cluster")

	// Create context with longer timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Determine which namespaces to search
	searchNamespaces := namespaces
	if allNs {
		// TODO: Implement getting all namespaces
		// For now, use default ones
		searchNamespaces = []string{"dev", "stage", "prod"}
		fmt.Println("📂 Searching in all namespaces...")
	} else {
		fmt.Printf("📂 Searching in namespaces: %v\n", searchNamespaces)
	}

	// Get Atlas applications
	apps, err := k8sClient.GetAtlasApps(ctx, searchNamespaces)
	if err != nil {
		return fmt.Errorf("failed to get Atlas applications: %w", err)
	}

	fmt.Printf("\n✅ Found %d Atlas application(s)\n\n", len(apps))

	// Display results in table format
	formatter.PrintAtlasAppsTable(apps)

	return nil
}
