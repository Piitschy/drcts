package cmd

import (
	"fmt"

	"github.com/Piitschy/drcts/internal/dialogs"
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
