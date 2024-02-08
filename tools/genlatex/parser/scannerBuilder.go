package parser

type ScannerBuilder interface {
	// Handle defines a Handler to invoke for a specific html.ElementNode type.
	// If a handler already exists then this Handler will be chained with it.
	Handle(string, Handler) ScannerBuilder

	// Default handler to use if an html.ElementNode is not matched.
	// If a handler already exists then this Handler will be chained with it.
	Default(Handler) ScannerBuilder

	// Comment Handler to invoke for html.CommentNode's
	// If a handler already exists then this Handler will be chained with it.
	Comment(Handler) ScannerBuilder

	// DocType defines the Handler to invoke for html.DoctypeNode's
	// If a handler already exists then this Handler will be chained with it.
	DocType(Handler) ScannerBuilder

	// Text defines the Handler to invoke for html.TextNode's
	// If a handler already exists then this Handler will be chained with it.
	Text(Handler) ScannerBuilder

	// Handler returns the Handler for this Builder.
	//
	// Unlike most builders, this will always return the same Handler for every call.
	// This allows for any changes made to the builder to be reflected in the Handler.
	//
	// This is important here because we have to support recursion between Scanner's.
	// e.g. content contains tables but table cells contain content.
	Handler() Handler
}

func New() ScannerBuilder {
	return &builder{
		scanner: Scanner{
			handlers: make(map[string]Handler),
		},
	}
}

type builder struct {
	scanner Scanner
}

func (b *builder) Handler() Handler {
	s := &b.scanner
	return s.Handle
}

func (b *builder) Handle(n string, h Handler) ScannerBuilder {
	b.scanner.handlers[n] = b.scanner.handlers[n].Then(h)
	return b
}

func (b *builder) Default(h Handler) ScannerBuilder {
	b.scanner.defaultHandler = b.scanner.defaultHandler.Then(h)
	return b
}

func (b *builder) Text(h Handler) ScannerBuilder {
	b.scanner.text = b.scanner.text.Then(h)
	return b
}

func (b *builder) Comment(h Handler) ScannerBuilder {
	b.scanner.text = b.scanner.text.Then(h)
	return b
}

func (b *builder) DocType(h Handler) ScannerBuilder {
	b.scanner.text = b.scanner.text.Then(h)
	return b
}
