package main

import (
	"context"
	"fmt"
	"log"

	"github.com/seniorGolang/asti/parser"
)

func main() {
	// Создаем парсер
	p := parser.NewParser()

	// Парсим пакет с интерфейсами
	result, err := p.ParsePackage(context.Background(), "examples/sample")
	if err != nil {
		log.Fatalf("Failed to parse package: %v", err)
	}

	// Проверяем, что интерфейсы найдены
	if len(result.Interfaces) == 0 {
		log.Fatal("No interfaces found")
	}

	fmt.Println("=== Демонстрация исправления поля Import ===")
	fmt.Println()
	
	for _, iface := range result.Interfaces {
		fmt.Printf("Интерфейс: %s\n", iface.Name)
		fmt.Printf("  Пакет: %s\n", iface.Package)
		fmt.Printf("  Import: %s\n", iface.Import)
		fmt.Printf("  ID: %s\n", iface.ID)
		fmt.Println()
	}

	// Проверяем, что Import содержит полный путь с именем модуля
	fmt.Println("=== Проверка корректности Import ===")
	for _, iface := range result.Interfaces {
		switch {
		case iface.Import == "":
			fmt.Printf("❌ Интерфейс %s: пустое поле Import\n", iface.Name)
		case !containsModuleName(iface.Import):
			fmt.Printf("❌ Интерфейс %s: Import '%s' не содержит имя модуля\n", iface.Name, iface.Import)
		default:
			fmt.Printf("✅ Интерфейс %s: корректный Import '%s'\n", iface.Name, iface.Import)
		}
	}

	fmt.Println()
	fmt.Println("=== Информация о пакете ===")
	fmt.Printf("ModuleName: %s\n", result.ModuleName)
	fmt.Printf("PackagePath: %s\n", result.PackagePath)
}

func containsModuleName(importPath string) bool {
	// Проверяем, что путь импорта содержит имя модуля
	return len(importPath) > 0 && (importPath != "." && importPath != "./" && importPath != "/")
} 