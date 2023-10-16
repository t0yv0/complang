package repl

import (
	"fmt"

	"context"
	"github.com/chzyer/readline"
	cl "github.com/t0yv0/complang"
	"github.com/t0yv0/complang/expr"
	"github.com/t0yv0/complang/parser"
)

type ReadEvalPrintLoopOptions struct {
	HistoryFile        string
	InitialEnvironment map[string]cl.Value
	MaxCompletions     int
}

func ReadEvalPrintLoop(ctx context.Context, cfg ReadEvalPrintLoopOptions) error {
	maxCompletions := cfg.MaxCompletions
	if maxCompletions == 0 {
		maxCompletions = 16
	}
	env := cl.NewMutableEnv()
	for k, v := range cfg.InitialEnvironment {
		env.Bind(k, v)
	}
	return readlineREPL(ctx, cfg.HistoryFile, &complangInterpreter{
		env:            env,
		maxCompletions: maxCompletions,
	})
}

type complangInterpreter struct {
	env            cl.MutableEnv
	maxCompletions int
}

func (ci *complangInterpreter) ReadEvalPrint(ctx context.Context, command string) {
	stmt, err := parser.ParseStmt(command)
	if err != nil {
		fmt.Printf("ERROR invalid syntax: %v\n", err)
		return
	}
	if stmt != nil {
		expr.EvalStmt(ctx, ci.env, stmt)
	}
}

func (ci *complangInterpreter) ReadEvalComplete(
	ctx context.Context,
	partialCommand string,
) []readline.Candidate {
	query, err := parser.ParseQuery(partialCommand)
	if err != nil {
		return nil
	}
	if query == nil {
		return nil
	}

	f := fuzzyComplete(ci.maxCompletions)
	expr.EvalQuery(ctx, ci.env, query, f.complete)

	result := []readline.Candidate{}
	for _, c := range f.matches(query.QueryText()) {
		rc := readline.Candidate{
			NewLine: []rune(partialCommand[0:query.Offset()] + c),
			Display: []rune(c),
		}
		result = append(result, rc)
	}
	return result
}
