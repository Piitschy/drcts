package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/Piitschy/drcts/internal/dialogs"
	"github.com/Piitschy/drcts/internal/directus"
	"github.com/urfave/cli/v2"
)

func ApplyDiff(cCtx *cli.Context) error {
	dTarget, err := dialogs.DirectusInstance(cCtx, "target")
	if err != nil {
		return err
	}

	if cCtx.String("input") == "" {
		return fmt.Errorf("Input file is required")
	}

	if path.Ext(cCtx.String("input")) != ".json" {
		return fmt.Errorf("Only json diffs can be applied")
	}

	diffBytes, err := os.ReadFile(cCtx.String("input"))
	if err != nil {
		return err
	}

	var diff directus.Diff
	err = json.Unmarshal(diffBytes, &diff)
	if err != nil {
		return err
	}

	sure, err := dialogs.YesNo(cCtx, "Do you want to apply and overwrite the current schema")
	if err != nil {
		return err
	}
	if !sure {
		return nil
	}
	verbose(cCtx, "Applying diff to target instance...")
	err = dTarget.ApplyDiff(&diff)
	if err != nil {
		return err
	}
	verbose(cCtx, "Diff applied successfully")
	return nil
}
