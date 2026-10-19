package cmd

import (
	"os"

	jlox "github.com/eaudetcobello/jlox/pkg"
	"github.com/spf13/cobra"
)

var code string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a jlox file",
	Run: func(cmd *cobra.Command, args []string) {
		if code != "" {
			if errs := jlox.Run(code); errs != nil {
				jlox.PrintError(os.Stderr, errs)
				os.Exit(1)
			}

			return
		}

		if len(args) == 1 {
			if err := jlox.RunFile(args[0]); err != nil {
				jlox.PrintError(os.Stderr, err)
				os.Exit(1)
			}

			return
		}

		jlox.RunInteractive(os.Stdin, os.Stdout, os.Stderr)
	},
}

func init() {
	runCmd.Flags().StringVarP(&code, "code", "c", "", "jlox code to run directly")
	rootCmd.AddCommand(runCmd)
}
