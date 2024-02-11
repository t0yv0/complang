package repl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/peterh/liner"
	cl "github.com/t0yv0/complang"
	"github.com/t0yv0/complang/expr"
	"github.com/t0yv0/complang/parser"
)

type ReadEvalPrintLoopOptions struct {
	HistoryFile        string
	InitialEnvironment map[string]cl.Value
	MaxCompletions     int
}

func ReadEvalPrintLoop(ctx context.Context, cfg ReadEvalPrintLoopOptions) (finalError error) {
	maxCompletions := cfg.MaxCompletions

	if maxCompletions == 0 {
		maxCompletions = 16
	}

	env := cl.NewMutableEnv()

	for k, v := range cfg.InitialEnvironment {
		env.Bind(k, v)
	}

	rliner := liner.NewLiner()

	defer func() {
		if err := rliner.Close(); err != nil {
			finalError = errors.Join(finalError, err)
		}
	}()

	rliner.SetCtrlCAborts(true)
	rliner.SetTabCompletionStyle(liner.TabPrints)

	rliner.SetWordCompleter(func(
		line string, pos int,
	) (head string, completions []string, tail string) {
		head = line[0:pos]
		completions = []string{}
		tail = line[pos:]
		queryText := line[0:pos]
		query, err := parser.ParseQuery(queryText)
		if err != nil {
			return
		}
		f := fuzzyComplete(maxCompletions)
		expr.EvalQuery(ctx, env, query, f.complete)
		head = line[0:query.Offset()]
		completions = f.matches(query.QueryText())
		tail = ""
		return
	})

	if cfg.HistoryFile != "" {
		if f, err := os.Open(cfg.HistoryFile); err == nil {
			_, err := rliner.ReadHistory(f)
			if err != nil {
				return fmt.Errorf("Error reading history file: %v", err)
			}
			if err := f.Close(); err != nil {
				return fmt.Errorf("Error closing history file: %v", err)
			}
		}
		defer func() {
			var e error
			if f, err := os.Create(cfg.HistoryFile); err != nil {
				e = fmt.Errorf("Error creating history file: %w", err)
			} else {
				if _, err := rliner.WriteHistory(f); err != nil {
					e = fmt.Errorf("Error writing history file: %w", err)
				} else if err := f.Close(); err != nil {
					e = fmt.Errorf("Error finalizing history file: %w", err)
				}
			}
			finalError = errors.Join(finalError, e)
		}()
	}

	for {
		if command, err := rliner.Prompt("> "); err == nil {
			stmt, err := parser.ParseStmt(command)
			if err != nil {
				fmt.Printf("ERROR invalid syntax: %v\n", err)
				continue
			}
			if stmt != nil {
				expr.EvalStmt(ctx, env, stmt)
				rliner.AppendHistory(command)
			}
		} else if err == liner.ErrPromptAborted || err == io.EOF {
			fmt.Println("")
			return nil
		} else {
			return fmt.Errorf("Error reading line: %w", err)
		}
	}
}
