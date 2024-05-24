package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Piitschy/drctsdm/internal/dialogs"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "drctsdm",
		Usage: "Directus Data Model CLI for schema migration\nBasic usage: drctsdm --base-url <base-url> --base-token <base-token> --target-url <target-url> --target-token <target-token> migrate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "base-url",
				Aliases: []string{"bu"},
				Usage:   "URL of the base Directus instance",
				EnvVars: []string{"DIRECTUS_BASE_URL"},
			},
			&cli.StringFlag{
				Name:    "base-token",
				Aliases: []string{"bt"},
				Usage:   "Token of the base Directus instance",
				EnvVars: []string{"DIRECTUS_BASE_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "target-url",
				Aliases: []string{"tu"},
				Usage:   "URL of the target Directus instance",
				EnvVars: []string{"DIRECTUS_TARGET_URL"},
			},
			&cli.StringFlag{
				Name:    "target-token",
				Aliases: []string{"tt"},
				Usage:   "Token of the target Directus instance",
				EnvVars: []string{"DIRECTUS_TARGET_TOKEN"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "Migrate schema",
				Action:  MigrateSchema,
			},
		},
	}

	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func MigrateSchema(cCtx *cli.Context) error {

	dBase, err := dialogs.DirectusInstance(cCtx, "base")
	dTarget, err := dialogs.DirectusInstance(cCtx, "target")
	if err != nil {
		fmt.Println(err)
	}

	s, err := dBase.GetSnapshot()
	diff, err := dTarget.GetDiff(s, false)

	if diff == nil {
		fmt.Println("No diff found")
		return nil
	}

	err = dTarget.ApplyDiff(diff)

	if err != nil {
		fmt.Println(err)
	}
	return nil
}
