package repl

import (
	"context"
	"testing"

	"github.com/chzyer/readline"
	"github.com/stretchr/testify/assert"
	cl "github.com/t0yv0/complang"
)

func TestReadEvalComplete(t *testing.T) {
	ctx := context.Background()
	env := cl.NewMutableEnv()
	env.Bind("$obj", cl.MapValue(map[string]cl.Value{
		"fox":  cl.StringValue{Text: "FOX"},
		"fine": cl.StringValue{Text: "FINE"},
	}))
	env.Bind("$fun", cl.BoolValue{Bool: true})
	inter := &complangInterpreter{env: env}

	{
		completions := inter.ReadEvalComplete(ctx, "$obj f")
		assert.Contains(t, completions, readline.Candidate{
			Display: []rune("fox"),
			NewLine: []rune("$obj fox"),
		})
		assert.Contains(t, completions, readline.Candidate{
			Display: []rune("fine"),
			NewLine: []rune("$obj fine"),
		})
	}

	{
		completions := inter.ReadEvalComplete(ctx, "$obj $f")
		for i, c := range completions {
			t.Logf("%d %q %q", i, string(c.Display), string(c.NewLine))
		}
		assert.Contains(t, completions, readline.Candidate{
			Display: []rune("$fun"),
			NewLine: []rune("$obj $fun"),
		})
	}
}
