package main

import (
	"context"

	"github.com/alecthomas/kong"

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

	cli := kong.Parse(
		&Commands{},
		kong.Name("movething"),
		kong.Description("a simple, adjustable file moving and renaming utility"),
		kong.BindTo(ctx, (*context.Context)(nil)),
	)
	err = cli.Run(ctx)
	cli.FatalIfErrorf(err)
}
