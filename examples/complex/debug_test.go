package complextestcases

import (
	"context"
	"fmt"
	"testing"

	"github.com/seniorGolang/asti/parser"
)

// TestDebugAnnotations тестирует парсинг аннотаций
func TestDebugAnnotations(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем парсинг sample пакета
	packagePath := "../sample"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse sample package: %v", err)
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
		if len(typeInfo.Fields) > 0 {
			fmt.Printf("    Fields: %d\n", len(typeInfo.Fields))
			for _, field := range typeInfo.Fields {
				fmt.Printf("      Field: %s\n", field.Name)
				fmt.Printf("        Annotations: %+v\n", field.Annotations)
			}
		}
	}
}
