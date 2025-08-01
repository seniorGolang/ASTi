package main

import (
	"context"
	"fmt"
	"log"

	"github.com/seniorGolang/asti/parser"
)

func main() {
	// Создаем парсер с дефолтным префиксом аннотаций @asti
	p := parser.NewParser()
	fmt.Printf("Default annotation prefix: %s\n", p.GetAnnotationPrefix())

	// Создаем парсер с кастомным префиксом аннотаций @custom
	pCustom := parser.NewParser(parser.WithAnnotationPrefix("@custom"))
	fmt.Printf("Custom annotation prefix: %s\n", pCustom.GetAnnotationPrefix())

	// Парсим пакет examples
	ctx := context.Background()
	pkg, err := p.ParsePackage(ctx, "examples/testplugin")
	if err != nil {
		log.Fatalf("Failed to parse package: %v", err)
	}

	// Выводим информацию о модуле и пакете
	fmt.Printf("Module: %s\n", pkg.ModuleName)
	fmt.Printf("Package: %s\n", pkg.PackagePath)
	fmt.Printf("Found %d interfaces:\n\n", len(pkg.Interfaces))

	for _, iface := range pkg.Interfaces {
		fmt.Printf("Interface: %s\n", iface.Name)
		fmt.Printf("  ID: %s\n", iface.ID)
		fmt.Printf("  Annotations: %v\n", iface.Annotations)
		fmt.Printf("  Methods: %d\n", len(iface.Methods))

		for _, method := range iface.Methods {
			fmt.Printf("    - %s (ID: %s)\n", method.Name, method.ID)
			fmt.Printf("      Annotations: %v\n", method.Annotations)
			fmt.Printf("      Parameters: %d\n", len(method.Parameters))
			fmt.Printf("      Results: %d\n", len(method.Results))
		}
		fmt.Println()
	}

	// Выводим информацию о типах
	fmt.Printf("Found %d types:\n", len(pkg.Types))
	for typeName, typeInfo := range pkg.Types {
		fmt.Printf("  - %s (%s)\n", typeName, typeInfo.Kind.String())
		if len(typeInfo.Fields) > 0 {
			fmt.Printf("    Fields: %d\n", len(typeInfo.Fields))
		}
		if len(typeInfo.Annotations) > 0 {
			fmt.Printf("    Annotations: %v\n", typeInfo.Annotations)
		}
	}

	// Сериализуем в JSON
	jsonData, err := p.ToJSON(pkg)
	if err != nil {
		log.Fatalf("Failed to serialize to JSON: %v", err)
	}

	fmt.Printf("\nJSON Output:\n%s\n", string(jsonData))
}
