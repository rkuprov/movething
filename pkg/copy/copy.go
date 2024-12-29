package copycmd

import (
	"context"
	"os"
	"path/filepath"
	"regexp"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"movething/pkg/logging"
	"movething/pkg/process"
)

type Cmd struct {
	Pattern     string `arg:""`
	Root        string `alias:"in" help:"Root directory to search for files"`
	Destination string `alias:"dest" help:"Destination directory to copy files to"`
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
	defaultRoot := filepath.Join(pwd, "testdata")
	defaultDestination := filepath.Join(pwd, "testdata", "dest")
	if c.Root == "" {
		c.Root = defaultRoot
	}
	if c.Destination == "" {
		c.Destination = defaultDestination
	}

	exp, err := regexp.Compile(c.Pattern)
	if err != nil {
		return err
	}

	toCopy, err := process.Config(ctx, c.Root, exp)
	if err != nil {
		return err
	}
	if len(toCopy) == 0 {
		logging.Info(ctx, "no targets found")
		return nil
	}
	logging.Info(ctx, "targets found", zap.Any("targets", toCopy))

	logging.Info(ctx, "copying files")
	logging.Info(ctx, "source", zap.String("src", toCopy[0].MatchPath))
	logging.Info(ctx, "destination", zap.String("dest", c.Destination))

	for _, item := range toCopy {
		err := copyItem(ctx, item, c.Root, c.Destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyItem(ctx context.Context, item process.PatternMatch, from, to string) error {
	f, err := os.Stat(item.MatchPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return copyItem(ctx, item, from, to)
	}
	err = os.MkdirAll(to, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Rename(item.MatchPath, filepath.Join(to, item.Match))
	if err != nil {
		return err
	}
	return nil
}
