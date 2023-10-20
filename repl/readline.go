package repl

import (
	"context"
	"io"
	"log"

	"github.com/chzyer/readline"
)

type interpreter interface {
	ReadEvalPrint(ctx context.Context, command string)
	ReadEvalComplete(partialCommand string) []readline.Candidate
}

func readlineREPL(
	ctx context.Context,
	historyFile string,
	inter interpreter,
) (finalError error) {
	cfg := &readline.Config{
		Prompt:            "\033[31m»\033[0m ",
		HistoryFile:       historyFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		AutoComplete:      &completer{inter: inter},
		FuncFilterInputRune: func(r rune) (rune, bool) {
			switch r {
			// block CtrlZ feature
			case readline.CharCtrlZ:
				return r, false
			}
			return r, true
		},
	}
	l, err := readline.NewEx(cfg)
	if err != nil {
		return err
	}
	defer func() {
		closeError := l.Close()
		if closeError != nil && finalError != nil {
			finalError = closeError
		}
	}()
	l.CaptureExitSignal()
	log.SetOutput(l.Stderr())

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		inter.ReadEvalPrint(ctx, line)
	}

	return nil
}

type completer struct {
	inter interpreter
}

var _ readline.AutoCompleterWithCandidates = (*completer)(nil)

// This method is not in use, Complete is called instead.
func (vc *completer) Do([]rune, int) ([][]rune, int) {
	panic("completer.Do is not supposed to be called")
}

func (vc *completer) Complete(
	line []rune, pos int,
) []readline.Candidate {
	// Only complete at the end of lines for now.
	if pos != len(line) {
		return nil
	}
	l := string(line)
	return vc.inter.ReadEvalComplete(l)
}
