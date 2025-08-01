package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// FindModuleRoot находит корень модуля Go для указанного пути
func FindModuleRoot(packagePath string) (moduleRoot string, err error) {
	// Используем golang.org/x/tools/go/packages для получения информации о модуле
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedModule,
		Dir:  packagePath,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		// Fallback: ищем go.mod файл вручную
		currentPath := packagePath
		for {
			goModPath := filepath.Join(currentPath, "go.mod")
			if _, statErr := os.Stat(goModPath); statErr == nil {
				moduleRoot = currentPath
				return
			}
			
			parent := filepath.Dir(currentPath)
			if parent == currentPath {
				// Достигли корня файловой системы
				break
			}
			currentPath = parent
		}
		
		err = fmt.Errorf("module root not found")
		return
	}

	if len(pkgs) == 0 {
		err = fmt.Errorf("no packages found")
		return
	}

	// Получаем корень модуля из Module.Dir
	mod := pkgs[0].Module
	if mod != nil {
		moduleRoot = mod.Dir
		return
	}
	
	err = fmt.Errorf("module root not found")
	return
}

// ParseGoMod парсит go.mod файл и извлекает имя модуля
func ParseGoMod(goModPath string) (moduleName string, err error) {
	content, err := os.ReadFile(goModPath)
	if err != nil {
		err = fmt.Errorf("failed to read go.mod file: %w", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			// Извлекаем имя модуля
			moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			// Убираем комментарии, если есть
			if commentIndex := strings.Index(moduleName, "//"); commentIndex != -1 {
				moduleName = strings.TrimSpace(moduleName[:commentIndex])
			}
			return
		}
	}

	err = fmt.Errorf("module declaration not found in go.mod file")
	return
} 