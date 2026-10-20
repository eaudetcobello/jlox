.PHONY: build test fmt ast repl

# Builds the jlox binary into src/jlox. This is the default target.
build:
	go build .

test:
	go test ./...

fmt:
	go mod tidy
	go fmt ./...

# Regenerates pkg/expr.go from the node list in cmd/generate_ast.go.
ast:
	go run . generate-ast pkg/expr.go

repl: build
	./jlox run
