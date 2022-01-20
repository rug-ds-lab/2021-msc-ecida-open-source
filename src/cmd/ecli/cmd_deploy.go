package main

import (
	"ecida/pkg/deployment"
	"ecida/pkg/depsresolver"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmdDeploy := &cobra.Command{
		Use:   "deploy <pipeline> <root-package>",
		Short: "deploy an ECiDA module to Kubernetes",
		Long:  `Deploy an ECiDA module to Kubernetes using the Helm structure.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			pipelineName := args[0]
			packageFile := args[1]

			return deployCmd(pipelineName, packageFile)
		},
	}

	rootCmd.AddCommand(cmdDeploy)
}

func deployCmd(pipelineName string, rootChart string) error {

	deps, err := depsresolver.FindDependencies(rootChart)

	if err != nil {
		return fmt.Errorf("unresolved dependencies for %s: %w\n", pipelineName, err)
	}

	for _, dep := range deps {
		fmt.Printf("%+v\n", dep.Values)
	}

	// deploy everything in deps to kubernetes
	err = deployment.Deploy(pipelineName, deps)

	if err != nil {
		return fmt.Errorf("failed to deploy %s: %w\n", pipelineName, err)
	}

	return nil
}
