package repl

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/t0yv0/complang/expr"
	"github.com/t0yv0/complang/parser"
	"github.com/t0yv0/complang/value"
)

type ReadEvalPrintLoopOptions struct {
	HistoryFile        string
	InitialEnvironment map[value.Symbol]value.Value
	MaxCompletions     int
}

func ReadEvalPrintLoop(cfg ReadEvalPrintLoopOptions) error {
	maxCompletions := cfg.MaxCompletions
	if maxCompletions == 0 {
		maxCompletions = 16
	}
	return readlineREPL(cfg.HistoryFile, &complangInterpreter{
		env:            value.MapEnv{cfg.InitialEnvironment},
		maxCompletions: maxCompletions,
	})
}

type complangInterpreter struct {
	env            value.MapEnv
	maxCompletions int
}

func (ci *complangInterpreter) ReadEvalPrint(command string) {
	stmt, err := parser.ParseStmt(command)
	if err != nil {
		fmt.Printf("ERROR invalid syntax: %v\n", err)
		return
	}
	if stmt != nil {
		expr.EvalStmt(&ci.env, stmt)
	}
}

func (ci *complangInterpreter) ReadEvalComplete(partialCommand string) []readline.Candidate {
	query, err := parser.ParseQuery(partialCommand)
	if err != nil {
		return nil
	}
	if query == nil {
		return nil
	}
	completions := expr.EvalQuery(&ci.env, query)
	result := []readline.Candidate{}
	for i, c := range completions {
		if i >= ci.maxCompletions {
			break
		}
		rc := readline.Candidate{
			NewLine: []rune(partialCommand[0:query.Offset()] + c.String()),
			Display: []rune(c.String()),
		}
		result = append(result, rc)
	}
	return result
}
