package cmd

import (
	"os"

	jlox "github.com/eaudetcobello/jlox/pkg"
	"github.com/spf13/cobra"
)

var code string

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run a jlox file, code passed with -c, or an interactive prompt",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if code == "" && len(args) == 0 {
			jlox.RunInteractive(cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())

			return
		}

		source, err := readSource(args)
		if err != nil {
			jlox.PrintError(cmd.ErrOrStderr(), err)
			os.Exit(1)
		}

		if err := jlox.Run(cmd.OutOrStdout(), cmd.ErrOrStderr(), source); err != nil {
			jlox.PrintError(cmd.ErrOrStderr(), err)
			os.Exit(1)
		}
	},
}

// readSource returns the code passed with -c, or else the contents of the file named in args.
func readSource(args []string) (string, error) {
	if code != "" {
		return code, nil
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func init() {
	runCmd.Flags().StringVarP(&code, "code", "c", "", "jlox code to run directly")
	rootCmd.AddCommand(runCmd)
}
