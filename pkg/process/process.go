package process

import (
	"context"
	"os"
	"path/filepath"
	"regexp"

	"go.uber.org/zap"

	"movething/pkg/logging"
)

type Task struct {
	Name                 string
	SearchDirectory      string
	SearchPattern        string
	DestinationDirectory string
	RenamePattern        string
}

type Matched struct {
	exp             *regexp.Regexp
	Pattern         string
	Match           string
	MatchPath       string
	CapturedMatches map[string]string
	DestinationPath string
}

func Config(ctx context.Context, task Task) ([]Matched, error) {
	items, err := os.ReadDir(task.SearchDirectory)
	if err != nil {
		return nil, err
	}

	exp, err := regexp.Compile(task.SearchPattern)
	if err != nil {
		return nil, err
	}
	var out []Matched
	for _, item := range items {

		if !exp.MatchString(item.Name()) {
			continue
		}
		if item.IsDir() {
			logging.Debug(ctx, "found directory", zap.String("dir", item.Name()))
			scanned, err := scanDirectory(ctx, filepath.Join(task.SearchDirectory, item.Name()), filepath.Join(task.DestinationDirectory, item.Name()), exp)
			if err != nil {
				return nil, err
			}

			if len(scanned) != 0 {
				out = append(out, scanned...)
			}

			continue
		}

		prepped, err := prepareFile(ctx, task.SearchDirectory, task.DestinationDirectory, item, exp)
		if err != nil {
			return nil, err
		}
		out = append(out, prepped)
	}

	return out, err
}

var acceptAll = regexp.MustCompile("^.*$")

func scanDirectory(ctx context.Context, current, to string, exp *regexp.Regexp) ([]Matched, error) {
	logging.Debug(ctx, "scanning directory", zap.String("dir", current))
	out := make([]Matched, 0)
	items, err := os.ReadDir(current)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		var scanned []Matched
		switch {
		case item.IsDir():
			scanned, err = scanDirectory(ctx,
				filepath.Join(current, item.Name()),
				filepath.Join(to, item.Name()),
				exp)
			if err != nil {
				return nil, err
			}

		default:
			f, err := prepareFile(ctx, current, to, item, acceptAll)
			if err != nil {
				return nil, err
			}
			logging.Debug(ctx, "found file", zap.String("file", f.Match))
			scanned = append(scanned, f)
		}

		out = append(out, scanned...)
	}

	return out, nil
}

func prepareFile(_ context.Context, searchDir, destinationDir string, file os.DirEntry, exp *regexp.Regexp) (Matched, error) {
	var match Matched
	match.exp = exp
	match.CapturedMatches = make(map[string]string)
	match.Pattern = exp.String()
	match.Match = file.Name()
	match.MatchPath = filepath.Join(searchDir, file.Name())

	submatches := exp.FindStringSubmatch(file.Name())
	groupNames := exp.SubexpNames()

	if len(groupNames) != 1 {
		for i, submatch := range submatches {
			if groupNames[i] == "" {
				continue
			}
			match.CapturedMatches[groupNames[i]] = submatch
		}
	}

	match.DestinationPath = filepath.Join(destinationDir, match.Match)

	return match, nil
}
