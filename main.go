package main

import (
	"context"

	"github.com/alecthomas/kong"

	copycmd "movething/pkg/copy"
	"movething/pkg/logging"
)

type Commands struct {
	Copy copycmd.Cmd `cmd:"" help:"Command to execute" required:""`
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
