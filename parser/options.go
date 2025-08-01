package parser

type Option func(parser *Parser)

func WithAnnotationPrefix(prefix string) Option {
	return func(parser *Parser) {
		parser.annotationPrefix = prefix
	}
}
