package dbacl

type Ast struct {
	Operator string
	Operands []AstNode
}

func (a Ast) Type() string {
	return "Ast"
}

type AstNode interface {
	Type() string
}

func parse(ts []token, i uint) (AstNode, error) {
	for i < uint(len(ts)) {
		
	}

	return nil, nil
}
