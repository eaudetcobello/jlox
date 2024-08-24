package ast

import (
	"fmt"
	"os"
	"strings"
)

func GenerateAst(outputDir string, baseName string, types []string) error {
	var builder strings.Builder

	if _, err := os.Stat(outputDir); err == nil {
		err := os.RemoveAll(outputDir)
		if err != nil {
			return err
		}
	}

	builder.WriteString("//+build meta\n\n")
	builder.WriteString("package ast\n\n")
	builder.WriteString(fmt.Sprintf("type %s struct {\n", baseName))
	builder.WriteString(fmt.Sprint("}\n"))

	err := os.WriteFile(outputDir, []byte(builder.String()), 0644)
	if err != nil {
		return err
	}

	return nil
}
