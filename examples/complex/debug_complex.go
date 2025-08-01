package complextestcases

import (
	"context"
	"fmt"
	"testing"

	"github.com/seniorGolang/asti/parser"
)

// TestDebugComplexImports тестирует парсинг complex_imports пакета
func TestDebugComplexImports(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем парсинг complex_imports пакета
	packagePath := "./complex_imports"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse complex_imports package: %v", err)
	}

	fmt.Printf("Package: %s\n", result.PackagePath)
	fmt.Printf("Package Annotations: %+v\n", result.Annotations)
	fmt.Printf("Found %d interfaces:\n", len(result.Interfaces))

	for _, iface := range result.Interfaces {
		fmt.Printf("  Interface: %s\n", iface.Name)
		fmt.Printf("    Annotations: %+v\n", iface.Annotations)
		fmt.Printf("    Methods: %d\n", len(iface.Methods))

		for _, method := range iface.Methods {
			fmt.Printf("      Method: %s\n", method.Name)
			fmt.Printf("        Annotations: %+v\n", method.Annotations)
		}
	}

	fmt.Printf("Found %d types:\n", len(result.Types))
	for typeName, typeInfo := range result.Types {
		fmt.Printf("  Type: %s\n", typeName)
		fmt.Printf("    Annotations: %+v\n", typeInfo.Annotations)
	}
}
