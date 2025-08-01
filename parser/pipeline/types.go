package pipeline

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"sync"

	"github.com/seniorGolang/asti/parser/models"
)

type StageTypeCollection struct {
	annotationParser models.AnnotationParser
	packageInfo      *models.Package
	imports          map[string]string // alias -> full path
	importsMutex     sync.RWMutex
}

func NewStageTypeCollection(annotationParser models.AnnotationParser) (stage *StageTypeCollection) {

	stage = &StageTypeCollection{annotationParser: annotationParser}
	return
}

// Process выполняет сбор типов
func (s *StageTypeCollection) Process(ctx context.Context, data Data) (result Data, err error) {

	// Package должен быть уже инициализирован предыдущими этапами
	if data.Package == nil {
		err = fmt.Errorf("package data is required for type collection")
		return
	}

	// Создаем карту для хранения всех типов
	allTypes := make(map[string]models.TypeInfo)

	// Множество уже обработанных типов для предотвращения рекурсии
	processedTypes := make(map[string]bool)

	// Сохраняем информацию о пакете для использования в других функциях
	s.packageInfo = data.Package

	// Получаем абсолютный путь для поиска файлов
	actualPackagePath := data.Package.PackagePath
	if data.Annotations != nil {
		if absPathData, exists := data.Annotations["_absolutePackagePath"]; exists {
			if absPath, ok := absPathData["path"]; ok {
				actualPackagePath = absPath
			}
		}
	}

	// Собираем информацию об импортах
	s.imports = make(map[string]string)
	s.collectImports(actualPackagePath)

	// Сначала собираем все типы из файлов
	pattern := filepath.Join(actualPackagePath, "*.go")
	files, err := filepath.Glob(pattern)
	if err == nil {
		fset := token.NewFileSet()
		for _, filename := range files {
			astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
			if err == nil {
				types, err := s.extractFromFile(context.Background(), astFile, fset, filename, actualPackagePath)
				if err == nil {
					for key, typeInfo := range types {
						allTypes[key] = typeInfo
					}
				}
			}
		}
	}

	// Затем проходим по отфильтрованным интерфейсам и собираем используемые типы
	for _, iface := range data.Interfaces {
		for _, method := range iface.Methods {
			for _, param := range method.Parameters {
				s.collectTypeRecursively(param.Type, actualPackagePath, allTypes, processedTypes)
			}

			for _, result := range method.Results {
				s.collectTypeRecursively(result.Type, actualPackagePath, allTypes, processedTypes)
			}
		}
	}

	data.Types = allTypes
	result = data
	return
}

// collectTypeRecursively собирает тип и все его зависимости рекурсивно
func (s *StageTypeCollection) collectTypeRecursively(typeStr string, packagePath string, usedTypes map[string]models.TypeInfo, processedTypes map[string]bool) {
	baseType := s.getBaseType(typeStr)
	if s.isBasicType(baseType) {
		return
	}

	// Проверяем, не обрабатывали ли мы уже этот тип
	if processedTypes[baseType] {
		return
	}

	// Помечаем тип как обрабатываемый
	processedTypes[baseType] = true

	// Ищем тип в файлах пакета
	typeInfo, found := s.findTypeInPackage(baseType, packagePath)
	if !found {
		// Если тип не найден в пакете, возможно это импортированный тип
		// Создаем базовую информацию для него
		packageName := s.getTypePackage(typeStr)
		importPath := ""
		
		if packageName == "" {
			// Если пакет не указан, используем текущий пакет
			if s.packageInfo != nil && s.packageInfo.ModuleName != "" {
				if s.packageInfo.PackagePath != "" && s.packageInfo.PackagePath != "." {
					packageName = s.packageInfo.ModuleName + "/" + s.packageInfo.PackagePath
				} else {
					packageName = s.packageInfo.ModuleName
				}
			} else {
				packageName = s.getPackageNameFromPath(packagePath)
			}
		} else {
			// Если пакет указан, ищем его в импортах
			s.importsMutex.RLock()
			if path, exists := s.imports[packageName]; exists {
				importPath = path
			}
			s.importsMutex.RUnlock()
		}

		// Для импортированных типов используем короткое имя пакета
		shortPackageName := packageName
		if importPath != "" {
			// Извлекаем короткое имя пакета из пути импорта
			parts := strings.Split(importPath, "/")
			shortPackageName = parts[len(parts)-1]
		}
		
		typeInfo = models.TypeInfo{
			Name:    s.getTypeName(baseType),
			Package: shortPackageName,
			Import:  importPath,
			Kind:    models.TypeBasic,
		}
	}

	// Сохраняем тип с полным именем (package.TypeName)
	packageName := s.getTypePackage(typeStr)
	if packageName == "" {
		// Если пакет не указан, используем текущий пакет
		packageName = s.getPackageNameFromPath(packagePath)
		fullTypeName := packageName + "." + baseType
		usedTypes[fullTypeName] = typeInfo
	} else {
		// Если пакет уже указан в типе, используем его как есть
		usedTypes[typeStr] = typeInfo
	}

	// Если это структура, рекурсивно собираем типы её полей
	if typeInfo.Kind == models.TypeStruct {
		for _, field := range typeInfo.Fields {
			s.collectTypeRecursively(field.Type, packagePath, usedTypes, processedTypes)
		}
	}
}

// findTypeInPackage ищет тип в файлах пакета
func (s *StageTypeCollection) findTypeInPackage(typeName string, packagePath string) (typeInfo models.TypeInfo, found bool) {
	// Получаем все Go файлы в пакете
	pattern := filepath.Join(packagePath, "*.go")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	fset := token.NewFileSet()
	for _, filename := range files {
		astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			continue
		}

		// Ищем тип в этом файле
		types, err := s.extractFromFile(context.Background(), astFile, fset, filename, packagePath)
		if err != nil {
			continue
		}

		// Ищем тип по имени
		for key, info := range types {
			if s.getTypeName(key) == typeName {
				typeInfo = info
				found = true
				return
			}
			// Также ищем по полному ключу, если typeName содержит точку
			if key == typeName {
				typeInfo = info
				found = true
				return
			}
		}
	}

	return
}

func (s *StageTypeCollection) getBaseType(typeStr string) (baseType string) {

	typeStr = strings.TrimPrefix(typeStr, "*")
	typeStr = strings.TrimPrefix(typeStr, "[]")
	if strings.HasPrefix(typeStr, "[") {
		if idx := strings.Index(typeStr, "]"); idx != -1 {
			typeStr = typeStr[idx+1:]
		}
	}
	typeStr = strings.TrimPrefix(typeStr, "chan ")
	typeStr = strings.TrimPrefix(typeStr, "chan<- ")
	typeStr = strings.TrimPrefix(typeStr, "<-chan ")
	typeStr = strings.TrimPrefix(typeStr, "...")

	// Обрабатываем дженерик типы - извлекаем имя типа до параметров
	if strings.Contains(typeStr, "[") {
		if idx := strings.Index(typeStr, "["); idx != -1 {
			typeStr = typeStr[:idx]
		}
	}
	baseType = typeStr
	return
}

func (s *StageTypeCollection) isBasicType(typeStr string) (isBasic bool) {

	if !strings.Contains(typeStr, ".") {
		basicTypes := map[string]bool{
			"string":      true,
			"int":         true,
			"int8":        true,
			"int16":       true,
			"int32":       true,
			"int64":       true,
			"uint":        true,
			"uint8":       true,
			"uint16":      true,
			"uint32":      true,
			"uint64":      true,
			"float32":     true,
			"float64":     true,
			"bool":        true,
			"byte":        true,
			"rune":        true,
			"error":       true,
			"interface{}": true,
			"any":         true,
		}
		isBasic = basicTypes[typeStr]
		return
	}

	isBasic = false
	return
}

func (s *StageTypeCollection) hasAnnotation(commentText string) (hasAnnotation bool) {

	commentText = strings.TrimSpace(commentText)
	if !strings.HasPrefix(commentText, "//") {
		hasAnnotation = false
		return
	}
	content := strings.TrimSpace(strings.TrimPrefix(commentText, "//"))
	if !strings.HasPrefix(content, "@") {
		hasAnnotation = false
		return
	}
	prefix := strings.TrimPrefix(s.annotationParser.(*models.DefaultAnnotationParser).GetPrefix(), "@")
	hasAnnotation = strings.HasPrefix(content[1:], prefix)
	

	
	return
}

func (s *StageTypeCollection) getTypeName(fullType string) (typeName string) {

	if idx := strings.LastIndex(fullType, "."); idx != -1 {
		typeName = fullType[idx+1:]
	} else {
		typeName = fullType
	}
	return
}

func (s *StageTypeCollection) getTypePackage(fullType string) (packageName string) {

	if idx := strings.LastIndex(fullType, "."); idx != -1 {
		packageName = fullType[:idx]
	} else {
		packageName = ""
	}
	return
}

// getPackageNameFromPath извлекает имя пакета из пути
func (s *StageTypeCollection) getPackageNameFromPath(packagePath string) (packageName string) {
	// Извлекаем последнюю часть пути как имя пакета
	parts := strings.Split(packagePath, "/")
	if len(parts) > 0 {
		packageName = parts[len(parts)-1]
		// Заменяем подчеркивания на точки для правильного имени пакета
		packageName = strings.ReplaceAll(packageName, "_", "")
	}
	return
}

func (s *StageTypeCollection) extractFromFile(ctx context.Context, astFile *ast.File, fset *token.FileSet, filename string, packagePath string) (types map[string]models.TypeInfo, err error) {

	types = make(map[string]models.TypeInfo)
	for _, decl := range astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					typeName := typeSpec.Name.Name
					fullTypeName := typeName
					if astFile.Name != nil {
						fullTypeName = astFile.Name.Name + "." + typeName
					}
					pos := fset.Position(typeSpec.Pos())
					var relativePath string
					relativePath, err = filepath.Rel(packagePath, filename)
					if err != nil {
						relativePath = filename
					}

					// Формируем полное имя пакета для импорта
					// Для типов из текущего пакета это должно быть полное имя модуля + путь пакета
					fullPackagePath := astFile.Name.Name
					fullImportPath := ""
					if s.packageInfo != nil && s.packageInfo.ModuleName != "" {
						// Для типов из текущего пакета формируем полное имя для импорта
						if s.packageInfo.PackagePath != "" && s.packageInfo.PackagePath != "." {
							fullImportPath = s.packageInfo.ModuleName + "/" + s.packageInfo.PackagePath
						} else {
							fullImportPath = s.packageInfo.ModuleName
						}
					}

					typeInfo := models.TypeInfo{
						Name:    typeName,
						Package: fullPackagePath,
						Import:  fullImportPath,
						Position: models.Position{
							File:   relativePath,
							Line:   pos.Line,
							Column: pos.Column,
						},
					}
					switch t := typeSpec.Type.(type) {
					case *ast.StructType:
						typeInfo.Kind = models.TypeStruct
						typeInfo.Fields = s.extractFields(ctx, t, fset, filename, packagePath)
					case *ast.InterfaceType:
						typeInfo.Kind = models.TypeInterface
					case *ast.ArrayType:
						typeInfo.Kind = models.TypeArray
						typeInfo.Array = true
					case *ast.MapType:
						typeInfo.Kind = models.TypeMap
						typeInfo.Map = true
					case *ast.ChanType:
						typeInfo.Kind = models.TypeChannel
						typeInfo.Channel = true
					case *ast.FuncType:
						typeInfo.Kind = models.TypeFunction
						typeInfo.Function = true
					default:
						typeInfo.Kind = models.TypeBasic
					}
					if typeSpec.TypeParams != nil {
						typeInfo.GenericType = true
						typeInfo.Generic = &models.GenericInfo{
							TypeParams: make([]string, 0, len(typeSpec.TypeParams.List)),
						}
						for _, param := range typeSpec.TypeParams.List {
							for _, name := range param.Names {
								typeInfo.Generic.TypeParams = append(typeInfo.Generic.TypeParams, name.Name)
							}
						}
					}
					if genDecl.Doc != nil {
						for _, comment := range genDecl.Doc.List {
							if s.hasAnnotation(comment.Text) {
								var annotations models.Annotations
								annotations, err = s.annotationParser.Parse(ctx, comment.Text)
								if err == nil {
									if typeInfo.Annotations == nil {
										typeInfo.Annotations = make(models.Annotations)
									}
									for k, v := range annotations {
										typeInfo.Annotations[k] = v
									}
								}
							}
						}
					}
					types[fullTypeName] = typeInfo
				}
			}
		}
	}
	return
}

func (s *StageTypeCollection) extractFields(ctx context.Context, structType *ast.StructType, fset *token.FileSet, filename string, packagePath string) (fields []models.FieldInfo) {

	if structType.Fields == nil {
		return fields
	}
	for _, field := range structType.Fields.List {
		fieldType := s.typeToString(field.Type)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				pos := fset.Position(name.Pos())
				var relativePath string
				var err error
				relativePath, err = filepath.Rel(packagePath, filename)
				if err != nil {
					relativePath = filename
				}
				fieldInfo := models.FieldInfo{
					Name: name.Name,
					Type: fieldType,
					Position: models.Position{
						File:   relativePath,
						Line:   pos.Line,
						Column: pos.Column,
					},
				}
				if field.Doc != nil {
					for _, comment := range field.Doc.List {
						if s.hasAnnotation(comment.Text) {
							var annotations models.Annotations
							annotations, err = s.annotationParser.Parse(ctx, comment.Text)
							if err == nil {
								if fieldInfo.Annotations == nil {
									fieldInfo.Annotations = make(models.Annotations)
								}
								for k, v := range annotations {
									fieldInfo.Annotations[k] = v
								}
							}
						}
					}
				}
				s.analyzeFieldType(field.Type, &fieldInfo)
				if field.Tag != nil {
					fieldInfo.Tags = s.parseTags(field.Tag.Value)
				}
				fields = append(fields, fieldInfo)
			}
		} else {
			pos := fset.Position(field.Pos())
			var relativePath string
			var err error
			relativePath, err = filepath.Rel(packagePath, filename)
			if err != nil {
				relativePath = filename
			}
			fieldInfo := models.FieldInfo{
				Name:     fieldType,
				Type:     fieldType,
				Embedded: true,
				Position: models.Position{
					File:   relativePath,
					Line:   pos.Line,
					Column: pos.Column,
				},
			}
			s.analyzeFieldType(field.Type, &fieldInfo)
			fields = append(fields, fieldInfo)
		}
	}
	return
}

func (s *StageTypeCollection) typeToString(expr ast.Expr) (typeStr string) {

	switch t := expr.(type) {
	case *ast.Ident:
		// Если это базовый тип, возвращаем как есть
		if s.isBasicType(t.Name) {
			typeStr = t.Name
		} else {
			// Проверяем, есть ли импорт с таким именем
			s.importsMutex.RLock()
			importPath, hasImport := s.imports[t.Name]
			s.importsMutex.RUnlock()
			
			if hasImport {
				// Это тип из внешнего пакета
				// Извлекаем имя пакета из пути импорта
				parts := strings.Split(importPath, "/")
				packageName := parts[len(parts)-1]
				typeStr = packageName + "." + t.Name
			} else {
				// Это тип из текущего пакета
				currentPackageName := s.getPackageNameFromPath(s.packageInfo.PackagePath)
				typeStr = currentPackageName + "." + t.Name
			}
		}
	case *ast.StarExpr:
		typeStr = "*" + s.typeToString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			typeStr = "[]" + s.typeToString(t.Elt)
		} else {
			typeStr = "[" + s.typeToString(t.Len) + "]" + s.typeToString(t.Elt)
		}
	case *ast.SelectorExpr:
		typeStr = s.typeToString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		typeStr = "interface{}"
	case *ast.MapType:
		typeStr = "map[" + s.typeToString(t.Key) + "]" + s.typeToString(t.Value)
	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			typeStr = "chan<- " + s.typeToString(t.Value)
		case ast.RECV:
			typeStr = "<-chan " + s.typeToString(t.Value)
		default:
			typeStr = "chan " + s.typeToString(t.Value)
		}
	case *ast.IndexExpr:
		// Для дженерик типов возвращаем базовый тип, параметры будут в analyzeFieldType
		typeStr = s.typeToString(t.X)
	case *ast.IndexListExpr:
		// Для дженерик типов возвращаем базовый тип, параметры будут в analyzeFieldType
		typeStr = s.typeToString(t.X)
	default:
		typeStr = "unknown"
	}
	return
}

func (s *StageTypeCollection) analyzeFieldType(expr ast.Expr, field *models.FieldInfo) {

	switch t := expr.(type) {
	case *ast.StarExpr:
		field.Pointer = true
		s.analyzeFieldType(t.X, field)
	case *ast.ArrayType:
		if t.Len == nil {
			field.Slice = true
		} else {
			field.Array = true
		}
		s.analyzeFieldType(t.Elt, field)
	case *ast.MapType:
		field.Map = true
		s.analyzeFieldType(t.Value, field)
	case *ast.ChanType:
		field.Channel = true
		s.analyzeFieldType(t.Value, field)
	case *ast.IndexExpr:
		field.Generic = true
		s.analyzeFieldType(t.X, field)
	case *ast.IndexListExpr:
		field.Generic = true
		s.analyzeFieldType(t.X, field)
	}
}

func (s *StageTypeCollection) parseTags(tagValue string) (tags map[string]string) {

	tags = make(map[string]string)
	tagValue = strings.Trim(tagValue, "`")
	parts := strings.Fields(tagValue)
	for _, part := range parts {
		if idx := strings.Index(part, ":"); idx != -1 {
			key := part[:idx]
			value := strings.Trim(part[idx+1:], `"`)
			tags[key] = value
		}
	}
	return
}

// collectImports собирает информацию об импортах из всех файлов пакета
func (s *StageTypeCollection) collectImports(packagePath string) {
	// Получаем все Go файлы в пакете
	pattern := filepath.Join(packagePath, "*.go")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	fset := token.NewFileSet()
	for _, filename := range files {
		astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			continue
		}

		// Обрабатываем импорты
		for _, decl := range astFile.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				for _, spec := range genDecl.Specs {
					if importSpec, ok := spec.(*ast.ImportSpec); ok {
						importPath := strings.Trim(importSpec.Path.Value, `"`)
						var alias string
						
						if importSpec.Name != nil {
							// Импорт с алиасом: import alias "path"
							alias = importSpec.Name.Name
						} else {
							// Импорт без алиаса: import "path"
							// Извлекаем имя пакета из пути
							parts := strings.Split(importPath, "/")
							alias = parts[len(parts)-1]
						}
						
						s.importsMutex.Lock()
						s.imports[alias] = importPath
						s.importsMutex.Unlock()
					}
				}
			}
		}
	}
}


