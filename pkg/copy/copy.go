package copycmd

import (
	"context"
	"fmt"
	"regexp"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"movething/pkg/logging"
	"movething/pkg/process"
)

type Cmd struct {
	Pattern     string `arg:""`
	Destination string `alias:"dest"`
	LogLevel    string `alias:"log" default:"info" enum:"debug,info,error" help:"Log level to use"`
}

func (c *Cmd) Run(ctx context.Context) error {
	switch c.LogLevel {
	case "debug":
		logging.SetLevel(zapcore.DebugLevel)
	case "error":
		logging.SetLevel(zapcore.ErrorLevel)
	}

	exp, err := regexp.Compile(c.Pattern)
	if err != nil {
		return err
	}

	toCopy, err := process.Config(ctx, c.Destination, exp)
	if err != nil {
		return err
	}

	for i := range len(toCopy) {
		var named []string
		for k, v := range toCopy[i].CapturedMatches {
			named = append(named, fmt.Sprintf("%s:%s", k, v))
		}

		if len(named) > 0 {
			logging.Debug(ctx, "matched", zap.String("directory", toCopy[i].Match), zap.Strings("named", named))
			continue
		}

		logging.Debug(ctx, "matched", zap.String("directory", toCopy[i].Match))
	}

	return err
}
