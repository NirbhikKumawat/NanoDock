package main

import (
	"fmt"
	"nanodocker/internal/editor"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "nanodock [file]",
		Short: "Terminal based Dockerfile Editor",
		Long:  "A Terminal based Dockerfile Editor",
		//Args:  cobra.MinimumNArgs(1),
		Run: editor.RunGocui,
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
