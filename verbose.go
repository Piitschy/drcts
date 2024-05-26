package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func verbose(cCtx *cli.Context, s string, a ...any) {
	if cCtx.Bool("verbose") {
		fmt.Printf(s+"\n", a...)
	}
}
