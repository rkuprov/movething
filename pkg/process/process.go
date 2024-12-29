package process

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
)

type Task struct {
	Name                 string
	Pattern              string
	RootDirectory        string
	SearchPattern        string
	DestinationDirectory string
	RenamePattern        string
}

type PatternMatch struct {
	Pattern         string
	Match           string
	MatchPath       string
	CapturedMatches map[string]string
}

func Config(_ context.Context, folder string, exp *regexp.Regexp) ([]PatternMatch, error) {
	items, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var out []PatternMatch
	for _, item := range items {
		var match PatternMatch
		match.CapturedMatches = make(map[string]string)
		match.Pattern = exp.String()

		if !exp.MatchString(item.Name()) {
			continue
		}
		match.Match = item.Name()
		match.MatchPath = filepath.Join(folder, item.Name())

		submatches := exp.FindStringSubmatch(item.Name())
		groupNames := exp.SubexpNames()

		if len(groupNames) != 1 {
			for i, submatch := range submatches {
				if groupNames[i] == "" {
					continue
				}
				match.CapturedMatches[groupNames[i]] = submatch
			}
		}

		out = append(out, match)
	}

	return out, err
}
