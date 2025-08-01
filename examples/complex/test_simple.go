package complextestcases

import (
	"context"
	"fmt"
	"testing"

	"github.com/seniorGolang/asti/parser"
)

// TestSimpleParsing тестирует простой парсинг
func TestSimpleParsing(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем парсинг sample пакета
	packagePath := "../sample"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse sample package: %v", err)
	}

	fmt.Printf("Package: %s\n", result.PackagePath)
	fmt.Printf("Found %d interfaces:\n", len(result.Interfaces))

	for _, iface := range result.Interfaces {
		fmt.Printf("  - %s\n", iface.Name)
		fmt.Printf("    Annotations: %v\n", iface.Annotations)
		fmt.Printf("    Methods: %d\n", len(iface.Methods))
	}

	fmt.Printf("Found %d types:\n", len(result.Types))
	for typeName := range result.Types {
		fmt.Printf("  - %s\n", typeName)
	}

	// Проверяем, что найдены интерфейсы
	if len(result.Interfaces) == 0 {
		t.Fatal("No interfaces found in sample package")
	}

	// Проверяем, что найдены типы
	if len(result.Types) == 0 {
		t.Fatal("No types found in sample package")
	}
}
