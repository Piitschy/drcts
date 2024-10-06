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

func ApplySchema(cCtx *cli.Context) error {
	dTarget, err := dialogs.DirectusInstance(cCtx, "target")
	if err != nil {
		return err
	}

	if cCtx.String("input") == "" {
		return fmt.Errorf("Input file is required")
	}

	if path.Ext(cCtx.String("input")) != ".json" {
		return fmt.Errorf("Only json snapshots can be applied")
	}

	verbose(cCtx, "Reading snapshot from "+cCtx.String("input"))
	snapshotBytes, err := os.ReadFile(cCtx.String("input"))
	if err != nil {
		return err
	}
	var s directus.Snapshot
	err = json.Unmarshal(snapshotBytes, &s)
	if err != nil {
		return err
	}
	for _, collectionName := range cCtx.Args().Slice() {
		s.FilterCollections(collectionName)
	}
	diff, err := dTarget.GetDiff(&s, true)
	if diff == nil {
		verbose(cCtx, "No changes detected")
		return nil
	}

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
	verbose(cCtx, "Applying schema to target instance...")
	err = dTarget.ApplyDiff(diff)
	if err != nil {
		return err
	}
	verbose(cCtx, "Schema applied successfully")
	return nil
}
