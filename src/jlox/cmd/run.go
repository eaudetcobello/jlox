/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	jlox "github.com/eaudetcobello/jlox/pkg"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a jlox file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			if err := jlox.RunFile(args[0]); err != nil {
				fmt.Println(fmt.Errorf("could not run file: %s", err))
				os.Exit(1)
			}
			os.Exit(0)
		}

		if err := jlox.RunPrompt(os.Stdin, os.Stdout, os.Stderr); err != nil {
			fmt.Println(fmt.Errorf("could not run prompt: %s", err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
