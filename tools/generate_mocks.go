//go:build generate

package main

import (
	"os"
	"text/template"
)

type MockFunction struct {
	Name        string
	Description string
	Signature   string
	Body        string
}

var mockFunctions = []MockFunction{
	{
		Name:        "MockTrueCondition",
		Description: "// MockCondition is a test condition implementation",
		Signature:   "func MockTrueCondition(ctx context.Context, data map[string]any) (bool, error)",
		Body:        "return true, nil",
	},
	{
		Name:        "MockFalseCondition",
		Description: "",
		Signature:   "func MockFalseCondition(ctx context.Context, data map[string]any) (bool, error)",
		Body:        "return false, nil",
	},
	{
		Name:        "MockErrorCondition",
		Description: "",
		Signature:   "func MockErrorCondition(ctx context.Context, data map[string]any) (bool, error)",
		Body:        "return false, errors.New(\"condition error\")",
	},
	{
		Name:        "MockSlowCondition",
		Description: "",
		Signature:   "func MockSlowCondition(ctx context.Context, data map[string]any) (bool, error)",
		Body:        "time.Sleep(300 * time.Millisecond)\n\treturn true, nil",
	},
	{
		Name:        "MockGuardCondition",
		Description: "// MockGuardCondition is a test guard condition",
		Signature:   "func MockGuardCondition(ctx context.Context, data map[string]any) (bool, error)",
		Body:        "return true, nil",
	},
	{
		Name:        "MockFailingGuardCondition",
		Description: "// MockFailingGuardCondition is a test guard condition that fails",
		Signature:   "func MockFailingGuardCondition(ctx context.Context, data map[string]any) (bool, error)",
		Body:        "return false, nil",
	},
	{
		Name:        "MockNoOpAction",
		Description: "// MockAction is a test action implementation",
		Signature:   "func MockNoOpAction(ctx context.Context, data map[string]any) (map[string]any, error)",
		Body:        "return nil, nil",
	},
	{
		Name:        "MockErrorAction",
		Description: "",
		Signature:   "func MockErrorAction(ctx context.Context, data map[string]any) (map[string]any, error)",
		Body:        "return nil, errors.New(\"action error\")",
	},
	{
		Name:        "MockSlowAction",
		Description: "",
		Signature:   "func MockSlowAction(ctx context.Context, data map[string]any) (map[string]any, error)",
		Body:        "select {\n\tcase <-time.After(300 * time.Millisecond):\n\t\treturn nil, nil\n\tcase <-ctx.Done():\n\t\treturn nil, ctx.Err()\n\t}",
	},
	{
		Name:        "MockUpdateAction",
		Description: "",
		Signature:   "func MockUpdateAction(ctx context.Context, data map[string]any) (map[string]any, error)",
		Body:        "return map[string]any{\n\t\t\"updated\": true,\n\t}, nil",
	},
}

const mockTemplate = `package machina

import (
	"context"
	"errors"
	"time"
)

{{range .}}
{{if .Description}}{{.Description}}
{{end}}{{.Signature}} {
	{{.Body}}
}

{{end}}`

func main() {
	tmpl, err := template.New("mocks").Parse(mockTemplate)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("mocks_test.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, mockFunctions)
	if err != nil {
		panic(err)
	}
}


