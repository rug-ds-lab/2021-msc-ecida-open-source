package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ecida/pkg/modulegen"

	"github.com/spf13/cobra"
)

func init() {
	cmdInit := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialise a new ECiDA module",
		Long:  "Initialise a new ECiDA module in a directory, or in current directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, _ := os.Getwd()

			if len(args) == 1 {
				err := makeDirIfNotExists(args[0])

				if err != nil {
					return err
				}

				path := filepath.Join(cwd, args[0])

				return initCmd(args[0], path)
			} else {
				_, name := filepath.Split(cwd)

				return initCmd(name, cwd)
			}
		},
	}

	rootCmd.AddCommand(cmdInit)
}

func makeDirIfNotExists(name string) error {
	_, err := os.Stat(name)

	if !os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s already exists", name))
	}

	return os.Mkdir(name, os.ModePerm)
}

func initCmd(name string, dirname string) error {
	fmt.Printf("initialising new module %s in %s\n", name, dirname)

	return modulegen.GenerateModule(name, dirname)
}
