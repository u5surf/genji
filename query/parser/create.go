package parser

import (
	"github.com/asdine/genji/query"
	"github.com/asdine/genji/query/scanner"
)

// parseCreateStatement parses a create string and returns a query.Statement AST object.
// This function assumes the CREATE token has already been consumed.
func (p *Parser) parseCreateStatement() (query.CreateTableStmt, error) {
	var stmt query.CreateTableStmt

	// Parse "TABLE".
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.TABLE {
		return stmt, newParseError(scanner.Tokstr(tok, lit), []string{"TABLE"}, pos)
	}

	// Parse table name
	tableName, err := p.ParseIdent()
	if err != nil {
		return stmt, err
	}
	stmt = query.CreateTable(tableName)

	// Parse "IF"
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok != scanner.IF {
		p.Unscan()
		return stmt, nil
	}

	// Parse "NOT"
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.NOT {
		return stmt, newParseError(scanner.Tokstr(tok, lit), []string{"NOT", "EXISTS"}, pos)
	}

	// Parse "EXISTS"
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.EXISTS {
		return stmt, newParseError(scanner.Tokstr(tok, lit), []string{"EXISTS"}, pos)
	}

	stmt = stmt.IfNotExists()

	return stmt, nil
}
