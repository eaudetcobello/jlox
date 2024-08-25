package ast

import (
	"fmt"
	"os"
	"strings"
)

// Parses a string in the format:
// Binary : left Expr, operator Token, right Expr
// Into a valid go struct definition
func defineType(builder *strings.Builder, structName string, fields string) {
	builder.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	for _, field := range strings.Split(fields, ",") {
		parts := strings.Fields(field)

		if len(parts) == 2 {
			name, fieldType := parts[0], parts[1]
			builder.WriteString(fmt.Sprintf("\t%s %s\n", name, fieldType))
		}
	}

	builder.WriteString(fmt.Sprint("}\n\n"))
}

func GenerateAst(outputDir string, types []string) error {
	var builder strings.Builder

	if _, err := os.Stat(outputDir); err == nil {
		err := os.RemoveAll(outputDir)
		if err != nil {
			return err
		}
	}

	builder.WriteString("package ast\n\n")

	builder.WriteString("import . \"github.com/eaudetcobello/jlox/pkg\"\n\n")

	builder.WriteString("type Expr interface{}\n\n")

	for _, t := range types {
		structName := strings.Replace(strings.Split(t, ":")[0], " ", "", -1)
		fields := strings.Split(t, ":")[1]

		defineType(&builder, structName, fields)
	}

	err := os.WriteFile(outputDir, []byte(builder.String()), 0644)
	if err != nil {
		return err
	}

	return nil
}
