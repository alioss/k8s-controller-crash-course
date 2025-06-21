package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cluster status information",
	Long: `Display the current status of the Kubernetes cluster including:
- Cluster name and version
- Number of users
- Node count
- Overall health status`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a sample Kubernetes cluster instance
		k8s := Kubernetes{
			Name:    "k8s-demo-cluster",
			Version: "1.31",
			Users:   []string{"alex", "den", "sarah"},
			NodeNumber: func() int {
				return 10
			},
		}

		// Display status information
		fmt.Printf("🎯 Cluster Status Report\n")
		fmt.Printf("========================\n")
		fmt.Printf("📛 Cluster Name: %s\n", k8s.Name)
		fmt.Printf("🔖 Version: %s\n", k8s.Version)
		fmt.Printf("👥 Total Users: %d\n", len(k8s.Users))
		fmt.Printf("🖥️  Active Nodes: %d\n", k8s.NodeNumber())
		fmt.Printf("✅ Status: Healthy\n")
		fmt.Printf("========================\n")
		
		// Show users list
		fmt.Printf("👤 Users:\n")
		for i, user := range k8s.Users {
			fmt.Printf("   %d. %s\n", i+1, user)
		}
	},
}

func init() {
	// Register the status command with the root command
	rootCmd.AddCommand(statusCmd)
}
