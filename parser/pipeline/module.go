package pipeline

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/seniorGolang/asti/parser/models"
)

type StageModule struct{}

func NewStageModule() (stage *StageModule) {
	stage = &StageModule{}
	return
}

// Process извлекает информацию о модуле из go.mod файла
func (s *StageModule) Process(ctx context.Context, data Data) (result Data, err error) {
	if data.Package == nil {
		err = fmt.Errorf("package data is required for module extraction")
		return
	}

	// Ищем go.mod файл, начиная с директории пакета и поднимаясь вверх
	moduleRootPath, moduleName, findErr := s.findModuleInfo(data.Package.PackagePath)
	if findErr != nil {
		// Не считаем это критической ошибкой, просто логируем и продолжаем
		fmt.Printf("Warning: failed to find module info: %v\n", findErr)
		// Для пустых директорий или директорий без модуля продолжаем работу
		// с абсолютным путем
		if data.Annotations == nil {
			data.Annotations = make(map[string]models.Annotations)
		}
		data.Annotations["_absolutePackagePath"] = models.Annotations{
			"path": data.Package.PackagePath,
		}
	} else {
		// Преобразуем абсолютные пути в относительные от корня модуля
		packagePath, err := filepath.Rel(moduleRootPath, data.Package.PackagePath)
		if err != nil {
			// Если не удалось получить относительный путь, используем абсолютный
			packagePath = data.Package.PackagePath
		}

		// Сохраняем абсолютный путь для внутреннего использования
		absolutePackagePath := data.Package.PackagePath

		data.Package.ModuleName = moduleName
		data.Package.PackagePath = packagePath

		// Добавляем абсолютный путь в данные для использования другими этапами
		if data.Annotations == nil {
			data.Annotations = make(map[string]models.Annotations)
		}
		data.Annotations["_absolutePackagePath"] = models.Annotations{
			"path": absolutePackagePath,
		}
	}

	result = data
	return
}

// findModuleInfo ищет go.mod файл и извлекает информацию о модуле
func (s *StageModule) findModuleInfo(packagePath string) (moduleRootPath, moduleName string, err error) {
	// Используем общую функцию для поиска корня модуля
	moduleRootPath, err = FindModuleRoot(packagePath)
	if err != nil {
		return
	}
	
	// Парсим go.mod файл для получения имени модуля
	goModPath := filepath.Join(moduleRootPath, "go.mod")
	moduleName, err = ParseGoMod(goModPath)
	return
}


