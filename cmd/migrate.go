package cmd

import (
	"fmt"

	"github.com/Piitschy/drcts/internal/dialogs"
	"github.com/Piitschy/drcts/internal/directus"
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
	if err != nil {
		return err
	}
	verbose(cCtx, "Getting diff from target instance...")
	diff, err := dTarget.GetDiff(s, force)
	if err == directus.DiffErr400 && !force {
		if !cCtx.Bool("yes") {
			yn, _ := dialogs.YesNo(cCtx, "Version mismatch between base and target instance. Do you want to force the migration anyway?")
			if !yn {
				return fmt.Errorf("%s Try to use --force to ignore this", err.Error())
			}
		} else {
			verbose(cCtx, "There seems to be a version mismatch between base and target instance. Trying to force migration...")
		}
		force = true
		verbose(cCtx, "Getting diff from target instance...")
		diff, err = dTarget.GetDiff(s, force)
	}
	if err != nil {
		return err
	}

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
