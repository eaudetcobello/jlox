package ast

import (
	"fmt"
	"os"
	"strings"
)

// Parses a string in the format:
// Binary : left Expr, operator Token, right Expr
// Into a valid go struct definition
func defineType(builder *strings.Builder, structName string, fields string) error {
	var tempBuilder strings.Builder

	// use a tempbuilder so that we don't write to the real builder
	// until we know no error has occurred
	tempBuilder.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	for _, field := range strings.Split(fields, ",") {
		parts := strings.Fields(field)

		if len(parts) == 2 {
			name, fieldType := parts[0], parts[1]
			tempBuilder.WriteString(fmt.Sprintf("  %s %s\n", name, fieldType))
			continue
		}

		return fmt.Errorf("field %v should have exactly two parts: field name and field name", parts)
	}

	tempBuilder.WriteString(fmt.Sprint("}\n\n"))

	builder.WriteString(tempBuilder.String())

	return nil
}

func GenerateAst(outputDir string, types []string) error {
	var builder strings.Builder

	builder.WriteString("package ast\n\n")

	builder.WriteString("import . \"github.com/eaudetcobello/jlox/pkg\"\n\n")

	builder.WriteString("type Expr interface{}\n\n")

	for _, t := range types {
		structName := strings.Replace(strings.Split(t, ":")[0], " ", "", -1)
		fields := strings.Split(t, ":")[1]

		if err := defineType(&builder, structName, fields); err != nil {
			return err
		}
	}


	if outputDir != "" {
		if _, err := os.Stat(outputDir); err == nil {
			err := os.RemoveAll(outputDir)
			if err != nil {
				return err
			}
		}

		err := os.WriteFile(outputDir, []byte(builder.String()), 0644)
		if err != nil {
			return err
		}

		return nil
	}

	fmt.Print(builder.String())

	return nil
}
