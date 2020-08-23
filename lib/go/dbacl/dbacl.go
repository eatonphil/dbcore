package dbacl

type Policy struct {
	ast AstNode
}

func (p Policy) ToSQL() string {
	return ""
}

func New(policy string) (*Policy, error) {
	tokens, err := lex(policy)
	if err != nil {
		return nil, err
	}

	sql, err := parse(tokens, 0)
	if err != nil {
		return nil, err
	}

	return &Policy{sql}, nil
}
