package cmd

import (
	"fmt"
	"os"

	"github.com/eaudetcobello/jlox/pkg/ast"
	"github.com/spf13/cobra"
)

var generateASTCmd = &cobra.Command{
	Use:   "generate-ast",
	Short: "Generate AST",
	Run: func(cmd *cobra.Command, args []string) {
		var filename string

		if len(args) > 0 {
			filename = args[0]
		} else {
			filename = ""
		}

		if err := ast.GenerateAst(filename, []string{
			"Binary : left Expr, operator Token , right Expr",
			"Grouping : expression Expr",
			"Literal : value interface{}",
			"Unary : operator Token, right Expr",
		}); err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateASTCmd)
}
