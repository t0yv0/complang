package repl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
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
	re, err := newRepl(ctx, cfg.MaxCompletions, cfg.InitialEnvironment, cfg.HistoryFile)
	if err != nil {
		return err
	}
	defer func() { finalError = errors.Join(finalError, re.close()) }()
	for re.active() {
		if err := re.interact(ctx); err != nil {
			return err
		}
	}
	return nil
}

type repl struct {
	prefix         string // used for post-processing fuzzy completion
	env            cl.MutableEnv
	maxCompletions int
	rliner         *liner.State
	stopped        bool
	historyFile    string
}

func newRepl(
	ctx context.Context, maxCompletions int, initEnv map[string]cl.Value, historyFile string,
) (*repl, error) {
	if maxCompletions == 0 {
		maxCompletions = 16
	}
	env := cl.NewMutableEnv()
	for k, v := range initEnv {
		env.Bind(k, v)
	}
	rliner := liner.NewLiner()
	rliner.SetCtrlCAborts(true)
	rliner.SetTabCompletionStyle(liner.TabPrints)
	re := &repl{
		maxCompletions: maxCompletions,
		rliner:         rliner,
		env:            env,
		historyFile:    historyFile,
	}
	rliner.SetWordCompleter(re.newWordCompleter(ctx))
	if err := re.readHistory(); err != nil {
		return nil, err
	}
	return re, nil
}

func (re *repl) interact(ctx context.Context) error {
	command, err := re.rliner.Prompt(re.prompt())
	command = re.prefix + command
	switch {
	case err == nil && strings.HasSuffix(command, "**"):
		newCommand, err := re.fuzzyFind(ctx, command)
		if err != nil {
			return fmt.Errorf("Error during fuzzy find: %v\n", err)
		}
		re.rliner.AppendHistory(command)
		re.rliner.AppendHistory(newCommand)
		re.prefix = newCommand
		return nil
	case err == nil:
		stmt, err := parser.ParseStmt(command)
		if err != nil {
			fmt.Printf("Error invalid syntax: %v\n", err)
			return nil
		}
		if stmt == nil {
			return nil
		}
		expr.EvalStmt(ctx, re.env, stmt)
		re.rliner.AppendHistory(command)
		re.prefix = ""
		return nil
	case err == liner.ErrPromptAborted || err == io.EOF:
		fmt.Println("")
		re.stopped = true
		return nil
	default:
		return fmt.Errorf("Error reading line: %w", err)
	}
}

func (re *repl) newWordCompleter(ctx context.Context) liner.WordCompleter {
	return func(line string, pos int) (head string, completions []string, tail string) {
		line = re.prefix + line
		pos = len(re.prefix) + pos
		defer func() { head = strings.TrimPrefix(head, re.prefix) }()
		head = line[0:pos]
		completions = []string{}
		tail = line[pos:]
		query, err := parser.ParseQuery(line[0:pos])
		if err != nil {
			return
		}
		expr.EvalQuery(ctx, re.env, query, func(_, candidate string) bool {
			if len(completions) > re.maxCompletions {
				return false
			}
			if strings.HasPrefix(candidate, query.QueryText()) {
				completions = append(completions, candidate)
			}
			return true
		})
		head = line[0:query.Offset()]
		return
	}
}

func (re *repl) fuzzyFind(ctx context.Context, command string) (string, error) {
	query, err := parser.ParseQuery(strings.TrimSuffix(command, "**"))
	if err != nil {
		return "", err
	}
	candidates := []string{}
	expr.EvalQuery(ctx, re.env, query, func(_, candidate string) bool {
		candidates = append(candidates, candidate)
		return true
	})
	selected, err := fuzzyfinder.Find(candidates,
		func(i int) string { return candidates[i] },
		fuzzyfinder.WithQuery(strings.ReplaceAll(command[query.Offset():], "**", "")))
	if err != nil && err == fuzzyfinder.ErrAbort {
		return command, nil
	} else if err != nil {
		return "", nil
	}
	return fmt.Sprintf("%s %s", command[0:query.Offset()], candidates[selected]), nil
}

func (re *repl) readHistory() error {
	if re.historyFile == "" {
		return nil
	}
	if f, err := os.Open(re.historyFile); err == nil {
		_, err := re.rliner.ReadHistory(f)
		if err != nil {
			return fmt.Errorf("Error reading history file: %v", err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("Error closing history file: %v", err)
		}
	}
	return nil
}

func (re *repl) writeHistory() error {
	if re.historyFile == "" {
		return nil
	}
	if f, err := os.Create(re.historyFile); err != nil {
		return fmt.Errorf("Error creating history file: %w", err)
	} else {
		if _, err := re.rliner.WriteHistory(f); err != nil {
			return fmt.Errorf("Error writing history file: %w", err)
		} else if err := f.Close(); err != nil {
			return fmt.Errorf("Error finalizing history file: %w", err)
		}
	}
	return nil
}

func (re *repl) close() error {
	return errors.Join(re.rliner.Close(), re.writeHistory())
}

func (re *repl) active() bool {
	return !re.stopped
}

func (re *repl) prompt() string {
	return fmt.Sprintf("> %s", re.prefix)
}
