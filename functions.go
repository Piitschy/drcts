package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/Piitschy/drctsdm/internal/dialogs"
	"github.com/Piitschy/drctsdm/internal/directus"
	"github.com/urfave/cli/v2"
)

func Migrate(cCtx *cli.Context) error {
	dBase, err := dialogs.DirectusInstance(cCtx, "base")
	dTarget, err := dialogs.DirectusInstance(cCtx, "target")
	if err != nil {
		return err
	}

	force := cCtx.Bool("force")

	verbose(cCtx, "Getting snapshot from base instance...")
	s, err := dBase.GetSnapshot()
	verbose(cCtx, "Getting diff from target instance...")
	diff, err := dTarget.GetDiff(s, force)

	if diff == nil {
		verbose(cCtx, "No changes detected")
		return nil
	}

	sure, err := dialogs.YesNo(cCtx, "Do you want to apply the changes to the target instance")
	if err != nil {
		return err
	}
	if !sure {
		return nil
	}
	verbose(cCtx, "Applying diff to target instance...")
	err = dTarget.ApplyDiff(diff)
	if err != nil {
		fmt.Println(err)
	}
	verbose(cCtx, "Migration completed")
	return nil
}

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
