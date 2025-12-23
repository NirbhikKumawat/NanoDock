package cmd

import (
	"fmt"
	"nanodocker/internal/editor"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "nanodock [file]",
	Short: "Terminal based Dockerfile Editor",
	Long:  "A Terminal based Dockerfile Editor",
	//Args:  cobra.MinimumNArgs(1),
	Run: editor.RunGocui,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
