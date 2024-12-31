package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"

	"movething/pkg/cfg"
	"movething/pkg/logging"
	"movething/pkg/move"
)

type Commands struct {
	Copy move.Cmd `cmd:"" help:"Command to execute" required:""`
}

func main() {
	ctx := context.Background()
	done, err := logging.SetupLogging(ctx)
	if err != nil {
		panic(err)
	}
	defer done()

	if len(os.Args) == 1 {
		err = tryWithConfigs(ctx)
		if err != nil {
			panic(err)
		}
		return
	}

	cli := kong.Parse(
		&Commands{},
		kong.Name("movething"),
		kong.Description("a simple, adjustable file moving and renaming utility"),
		kong.BindTo(ctx, (*context.Context)(nil)),
	)
	err = cli.Run(ctx)
	cli.FatalIfErrorf(err)
}

func tryWithConfigs(ctx context.Context) error {
	configs, err := cfg.GetConfig(ctx)
	if err != nil {
		return err
	}

	for _, task := range configs {
		m := move.Cmd{
			Pattern:     task.SearchPattern,
			Root:        task.SearchDirectory,
			Destination: task.DestinationDirectory,
		}
		err = m.Run(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
