package repl

import (
	"sort"
	"strings"

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
		return sortedCompletions(query, targets, maxCompletions)
	}

	return fuzzyCompleter{
		complete: completion,
		matches:  matches,
	}
}

func sortedCompletions(query string, targets []string, maxCompletions int) []string {
	ranks := fuzzy.RankFindNormalizedFold(query, targets)
	var out []string
	sort.SliceStable(ranks, func(i, j int) bool {
		istr := targets[ranks[i].OriginalIndex]
		jstr := targets[ranks[j].OriginalIndex]
		if strings.HasPrefix(istr, query) && !strings.HasPrefix(jstr, query) {
			return true
		}
		if strings.Contains(istr, query) && !strings.Contains(jstr, query) {
			return true
		}
		return ranks[i].Distance < ranks[j].Distance
	})
	for _, r := range ranks {
		out = append(out, targets[r.OriginalIndex])
		if len(out) > maxCompletions {
			break
		}
	}
	return out
}
