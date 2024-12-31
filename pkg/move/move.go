package move

import (
	"context"
	"os"
	"path/filepath"
	"slices"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"movething/pkg/logging"
	"movething/pkg/process"
)

type Cmd struct {
	Pattern     string `arg:""`
	Root        string `alias:"in" help:"Root directory to search for files"`
	Destination string `alias:"to" help:"Destination directory to move files to"`
	LogLevel    string `alias:"log" default:"info" enum:"debug,info,error" help:"Log level to use"`
}

func (c *Cmd) Run(ctx context.Context) error {
	switch c.LogLevel {
	case "debug":
		logging.SetLevel(zapcore.DebugLevel)
	case "error":
		logging.SetLevel(zapcore.ErrorLevel)
	}
	var pwd, _ = os.Getwd()
	defaultRoot := filepath.Join(pwd, "testdata", "from")
	defaultDestination := filepath.Join(pwd, "testdata", "to")
	if c.Root == "" {
		c.Root = defaultRoot
	}
	if c.Destination == "" {
		c.Destination = defaultDestination
	}

	toCopy, err := process.Config(ctx, process.Task{
		Name:                 "",
		SearchDirectory:      c.Root,
		SearchPattern:        c.Pattern,
		DestinationDirectory: c.Destination,
		RenamePattern:        "",
	})
	if err != nil {
		return err
	}
	if len(toCopy) == 0 {
		logging.Info(ctx, "no targets found")
		return nil
	}

	logging.Info(ctx, "copying files")
	toDelete := make([]string, 0)
	for _, item := range toCopy {
		logging.Debug(ctx, "moving file",
			zap.String("file", item.Match),
			zap.String("from", item.MatchPath),
			zap.String("to", item.DestinationPath))
		err = os.MkdirAll(filepath.Dir(item.DestinationPath), 0755)
		if err != nil {
			return err
		}

		err = os.Rename(item.MatchPath, item.DestinationPath)
		if err != nil {
			return err
		}

		if filepath.Dir(item.MatchPath) != c.Root && !slices.Contains(toDelete, filepath.Dir(item.MatchPath)) {
			logging.Debug(ctx, "queued to remove", zap.String("directory", filepath.Dir(item.MatchPath)))
			defer os.RemoveAll(filepath.Dir(item.MatchPath))
			toDelete = append(toDelete, filepath.Dir(item.MatchPath))
		}
	}

	return nil
}
