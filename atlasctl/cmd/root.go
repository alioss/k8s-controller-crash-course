package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "atlasctl",
	Short: "Atlas deployment management tool",
	Long: `AtlasCtl is a CLI tool for managing Atlas application deployments 
across different environments (dev, stage, prod) in Kubernetes.

It helps you monitor versions, migration IDs, and promote deployments 
between environments.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
}
