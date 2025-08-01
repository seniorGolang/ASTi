package models

import (
	"context"
	"strings"
)

type Annotations map[string]string

type AnnotationParser interface {
	Parse(ctx context.Context, text string) (annotations Annotations, err error)
}

type DefaultAnnotationParser struct {
	prefix string
}

func NewAnnotationParser(prefix string) (parser *DefaultAnnotationParser) {

	if prefix == "" {
		prefix = "@asti"
	}

	parser = &DefaultAnnotationParser{prefix: prefix}
	return
}

func (p *DefaultAnnotationParser) GetPrefix() (prefix string) {

	prefix = p.prefix
	return
}

// Parse парсит аннотации из текста комментария
func (p *DefaultAnnotationParser) Parse(_ context.Context, text string) (annotations Annotations, err error) {

	annotations = make(Annotations)

	text = strings.TrimSpace(text)

	if strings.HasPrefix(text, "//") {
		text = strings.TrimSpace(strings.TrimPrefix(text, "//"))
	}

	if !strings.HasPrefix(text, p.prefix) {
		return
	}

	content := strings.TrimSpace(strings.TrimPrefix(text, p.prefix))

	var currentKey string
	var currentValue strings.Builder
	inQuotes := false

	tokens := p.tokenize(content)

	for _, token := range tokens {
		switch {
		case strings.Contains(token, "=") && !inQuotes:
			if currentKey != "" {
				annotations[currentKey] = strings.TrimSpace(currentValue.String())
				currentValue.Reset()
			}
			if key, value, found := strings.Cut(token, "="); found {
				currentKey = strings.TrimSpace(key)
				value = strings.TrimSpace(value)

				if strings.HasPrefix(value, `"`) {
					if strings.HasSuffix(value, `"`) {
						annotations[currentKey] = strings.Trim(value, `"`)
						currentKey = ""
					} else {
						currentValue.WriteString(value[1:])
						inQuotes = true
					}
				} else {
					annotations[currentKey] = value
					currentKey = ""
				}
			}
		case inQuotes:
			if strings.HasSuffix(token, `"`) {
				currentValue.WriteString(" " + strings.TrimSuffix(token, `"`))
				annotations[currentKey] = strings.TrimSpace(currentValue.String())
				currentKey = ""
				currentValue.Reset()
				inQuotes = false
			} else {
				currentValue.WriteString(" " + token)
			}
		case currentKey != "":
			currentValue.WriteString(" " + token)
		default:
			// Обработка короткой записи булевых значений (только ключ без значения)
			// Аннотация вида @asti key интерпретируется как @asti key=true
			annotations[strings.TrimSpace(token)] = "true"
		}
	}

	if currentKey != "" {
		annotations[currentKey] = strings.TrimSpace(currentValue.String())
	}

	return
}

// tokenize разбивает строку аннотации на токены
func (p *DefaultAnnotationParser) tokenize(content string) (tokens []string) {

	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(content); i++ {
		char := content[i]

		switch {
		case char == '"':
			inQuotes = !inQuotes
			current.WriteByte(char)
		case char == ' ' && !inQuotes:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return
}
