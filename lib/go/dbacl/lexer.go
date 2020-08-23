package dbacl

type kind uint

const (
	identKind kind = iota
	symbolKind
)

type token struct {
	loc uint
	value string
	kind kind
}

func lex(policy string) ([]token, error) {
	return nil, nil
}
