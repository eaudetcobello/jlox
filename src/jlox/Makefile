.PHONY: go.fmt
go.fmt:
	go mod tidy
	go fmt ./...

.PHONY: ast
ast:
	go run main.go generate-ast pkg/expr.go

.PHONY: test
test:
	go test ./...

build:
	go build .

run-interactive: build
	go run . run
