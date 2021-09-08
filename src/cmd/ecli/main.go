package main

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ecli",
		Short: "CLI tool to interact with ECiDA",
		Long:  "CLI tool to package ECiDA modules, interact with a cluster, and deploy said modules",
	}
)

func main() {
    // The init() funciton in each of the cmd_ files in this package add
    // commands to the cobra rootCmd.
	rootCmd.Execute()
}
