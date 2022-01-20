package main

import "os"

func DirOrWorkdir(dir string) string {
	if dir != "" {
		return dir
	}

	workdir, _ := os.Getwd()

	return workdir
}
