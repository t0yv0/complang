package expr

type Expr interface {
	exprMarker()
}

type RefExpr struct {
	exprMarkerImpl
	Ref string
	// Character-based offset of the start of the symbol appearance in the source code.
	Offset int
}

var _ Expr = (*RefExpr)(nil)

type MessageExpr struct {
	exprMarkerImpl
	Receiver Expr
	Message  Expr
}

var _ Expr = (*MessageExpr)(nil)

type SymbolExpr struct {
	exprMarkerImpl
	Symbol string
	// Character-based offset of the start of the symbol appearance in the source code.
	Offset int
}

var _ Expr = (*SymbolExpr)(nil)

type NullExpr struct {
	exprMarkerImpl
}

var _ Expr = (*NullExpr)(nil)

type BoolExpr struct {
	exprMarkerImpl
	Bool bool
}

var _ Expr = (*BoolExpr)(nil)

type StringExpr struct {
	exprMarkerImpl
	String string
}

var _ Expr = (*StringExpr)(nil)

type LambdaBlockExpr struct {
	exprMarkerImpl
	Symbols []string
	Body    Expr
}

var _ Expr = (*LambdaBlockExpr)(nil)

type exprMarkerImpl struct{}

func (*exprMarkerImpl) exprMarker() {}

type Stmt interface {
	stmtMarker()
}

type ExprStmt struct {
	stmtMarkerImpl
	Expr Expr
}

var _ Stmt = (*ExprStmt)(nil)

type AssignStmt struct {
	stmtMarkerImpl
	Ref  string
	Expr Expr
}

var _ Stmt = (*AssignStmt)(nil)

type stmtMarkerImpl struct{}

func (*stmtMarkerImpl) stmtMarker() {}

type Query interface {
	queryMarker()
	// Source-based offset (starting at 0) indicating the position of the lexeme being completed.
	Offset() int

	QueryText() string
}

type SymbolQuery struct {
	queryMarkerImpl
	Expr         Expr
	Symbol       string
	SymbolOffset int
}

var _ Query = (*SymbolQuery)(nil)

func (sq *SymbolQuery) Offset() int {
	return sq.SymbolOffset
}

func (sq *SymbolQuery) QueryText() string {
	return sq.Symbol
}

type RefQuery struct {
	queryMarkerImpl
	Ref       string
	RefOffset int
}

var _ Query = (*RefQuery)(nil)

func (rq *RefQuery) Offset() int {
	return rq.RefOffset
}

func (rq *RefQuery) QueryText() string {
	return rq.Ref
}

type queryMarkerImpl struct{}

func (*queryMarkerImpl) queryMarker() {}
