package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	cmdPublish := &cobra.Command{
		Use:   "publish",
		Short: "publish a package",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	rootCmd.AddCommand(cmdPublish)
}

func publishCmd(file string) error {
	dir, filename := filepath.Split(file)

	addCmd := exec.Command("git", "add", filename)
	addCmd.Dir = dir
	err := addCmd.Run()

	if err != nil {
		return err
	}

	commitCmd := exec.Command("git", "commit", "-m", fmt.Sprintf("publish %s", filename))
	commitCmd.Dir = dir
	err = commitCmd.Run()

	if err != nil {
		return err
	}

	pushCmd := exec.Command("git", "push")
	pushCmd.Dir = dir

	err = pushCmd.Run()

	return err
}
