package repl

import (
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type fuzzyCompleter struct {
	complete func(query string, target string) bool
	matches  func(query string) []string
}

func fuzzyComplete(maxCompletions int) fuzzyCompleter {
	targets := []string{}

	completion := func(query string, target string) bool {
		targets = append(targets, target)
		return true
	}

	matches := func(query string) []string {
		var out []string
		ranks := fuzzy.RankFindNormalizedFold(query, targets)
		sort.Slice(ranks, func(i, j int) bool {
			return ranks[i].Distance < ranks[j].Distance
		})
		for _, r := range ranks {
			out = append(out, targets[r.OriginalIndex])
		}
		return out
	}

	return fuzzyCompleter{
		complete: completion,
		matches:  matches,
	}
}
