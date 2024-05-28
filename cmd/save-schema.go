package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/Piitschy/drcts/internal/dialogs"
	"github.com/urfave/cli/v2"
)

func SaveSchema(cCtx *cli.Context) error {
	dBase, err := dialogs.DirectusInstance(cCtx, "base")
	if err != nil {
		return err
	}

	fileExt := path.Ext(cCtx.String("output"))

	if fileExt != "" && cCtx.String("format") != "" {
		return fmt.Errorf("Output file extension and format flag are mutually exclusive")
	}
	if fileExt == "" && cCtx.String("format") == "" {
		return fmt.Errorf("Output file extension or format flag is required")
	}

	if fileExt == "" && cCtx.String("format") != "" {
		fileExt = "." + cCtx.String("format")
	}

	verbose(cCtx, "Getting schema from base instance...")
	b, err := dBase.GetRawSnapshot(fileExt[1:])

	if err != nil {
		return err
	}

	if cCtx.String("output") == "" {
		os.Stdout.Write(b)
		return nil
	}

	err = os.WriteFile(cCtx.String("output"), b, 0644)
	if err != nil {
		return err
	}
	verbose(cCtx, "Schema written successfully to "+cCtx.String("output"))

	return nil
}
