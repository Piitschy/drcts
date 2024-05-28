package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/Piitschy/drcts/internal/dialogs"
	"github.com/Piitschy/drcts/internal/directus"
	"github.com/urfave/cli/v2"
)

func SaveDiff(cCtx *cli.Context) error {
	if cCtx.String("input") != "" {
		if path.Ext(cCtx.String("input")) != ".json" {
			return fmt.Errorf("Only json diffs can be applied")
		}
	}

	var dBase *directus.Directus
	if cCtx.String("input") == "" {
		var err error
		dBase, err = dialogs.DirectusInstance(cCtx, "base")
		if err != nil {
			return err
		}
	}

	dTarget, err := dialogs.DirectusInstance(cCtx, "target")
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

	if fileExt != ".json" {
		fmt.Println("Attention!! Only json diffs can be applyed")
	}

	diffBytes := []byte{}
	if cCtx.String("input") == "" {
		s, err := dBase.GetSnapshot()
		if err != nil {
			return err
		}
		diffBytes, err = dTarget.GetRawDiff(s, fileExt[1:], cCtx.Bool("force"))
	} else {
		diffBytes, err = os.ReadFile(cCtx.String("input"))
	}

	if err != nil {
		return err
	}

	if cCtx.String("output") == "" {
		os.Stdout.Write(diffBytes)
		return nil
	}
	os.WriteFile(cCtx.String("output"), diffBytes, 0644)
	verbose(cCtx, "Diff written successfully to "+cCtx.String("output"))

	return nil
}
