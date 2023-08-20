package repl

import (
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

type interpreter interface {
	ReadEvalPrint(command string)
	ReadEvalComplete(partialCommand string) (string, []string)
}

func readlineREPL(historyFile string, inter interpreter) (finalError error) {
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
		inter.ReadEvalPrint(line)
	}

	return nil
}

type completer struct {
	inter interpreter
}

func (vc *completer) Do(line []rune, pos int) ([][]rune, int) {
	// Only complete at the end of lines for now.
	if pos != len(line) {
		return nil, len(line)
	}
	l := runesToStr(line)
	prefix, completions := vc.inter.ReadEvalComplete(l)
	lastToken := l[len(prefix):]
	var out [][]rune
	for _, x := range completions {
		if strings.HasPrefix(x, lastToken) {
			out = append(out, strToRunes(strings.TrimPrefix(x, lastToken)))
		}
	}
	return out, len(lastToken)
}
