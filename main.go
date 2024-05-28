package main

import (
	"log"
	"os"

	"github.com/Piitschy/drcts/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app.Suggest = true
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var app = &cli.App{
	Name:  "drcts",
	Usage: "Directus Data Model CLI for schema migration\nBasic usage: drcts --base-url <base-url> --base-token <base-token> --target-url <target-url> --target-token <target-token> migrate",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "base-url",
			Aliases: []string{"bu"},
			Usage:   "URL of the base Directus instance",
			EnvVars: []string{"DRCTS_BASE_URL"},
		},
		&cli.StringFlag{
			Name:    "base-token",
			Aliases: []string{"bt"},
			Usage:   "Token of the base Directus instance",
			EnvVars: []string{"DRCTS_BASE_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "target-url",
			Aliases: []string{"tu"},
			Usage:   "URL of the target Directus instance",
			EnvVars: []string{"DRCTS_TARGET_URL"},
		},
		&cli.StringFlag{
			Name:    "target-token",
			Aliases: []string{"tt"},
			Usage:   "Token of the target Directus instance",
			EnvVars: []string{"DRCTS_TARGET_TOKEN"},
		},
	},
	Commands: []*cli.Command{
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "Migrate schema from base to target instance",
			Action:  cmd.Migrate,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "force",
					Aliases:     []string{"f"},
					DefaultText: "false",
					Usage:       "Ignore version and db differences between instances",
					EnvVars:     []string{"DRCTS_FORCE"},
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"v"},
					DefaultText: "false",
					Usage:       "Verbose",
				},
				&cli.BoolFlag{
					Name:        "yes",
					Aliases:     []string{"y"},
					DefaultText: "false",
					Usage:       "Don't ask for confirm",
				},
			},
		},
		{
			Name:    "save-schema",
			Aliases: []string{"ss", "save"},
			Usage:   "Save snapshot of the base instance to output file (-o) or stdout",
			Action:  cmd.SaveSchema,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:      "output",
					Aliases:   []string{"o"},
					Usage:     "Output file",
					EnvVars:   []string{"DRCTS_SCHEMA_FILE"},
					TakesFile: true,
				},
				&cli.StringFlag{
					Name:    "format",
					Aliases: []string{"F"},
					Usage:   "Output format (json, yaml, xml, csv)",
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"v"},
					DefaultText: "false",
					Usage:       "Verbose",
				},
			},
		},
		{
			Name:    "save-diff",
			Aliases: []string{"sd", "diff"},
			Usage:   "Save diff of the base and target instances to output file (-o) or stdout",
			Action:  cmd.SaveDiff,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:      "input",
					Aliases:   []string{"i"},
					Usage:     "Give a schema file to compare with the target instance",
					EnvVars:   []string{"DRCTS_SCHEMA_FILE"},
					TakesFile: true,
				},
				&cli.StringFlag{
					Name:      "output",
					Aliases:   []string{"o"},
					Usage:     "Output file",
					EnvVars:   []string{"DRCTS_DIFF_FILE"},
					TakesFile: true,
				},
				&cli.StringFlag{
					Name:    "format",
					Aliases: []string{"F"},
					Usage:   "Output format (json, yaml, xml, csv)",
				},
				&cli.BoolFlag{
					Name:        "force",
					Aliases:     []string{"f"},
					DefaultText: "false",
					Usage:       "Ignore version and db differences between instances",
					EnvVars:     []string{"DRCTS_FORCE"},
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"v"},
					DefaultText: "false",
					Usage:       "Verbose",
				},
			},
		},
		{
			Name:    "apply-diff",
			Aliases: []string{"ad"},
			Usage:   "Apply diff to target instance",
			Action:  cmd.ApplyDiff,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:      "input",
					Aliases:   []string{"i"},
					Usage:     "Input file",
					EnvVars:   []string{"DRCTS_DIFF_FILE"},
					TakesFile: true,
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"v"},
					DefaultText: "false",
					Usage:       "Verbose",
				},
				&cli.BoolFlag{
					Name:        "yes",
					Aliases:     []string{"y"},
					DefaultText: "false",
					Usage:       "Don't ask for confirm",
				},
			},
		},
		{
			Name:    "apply-schema",
			Aliases: []string{"as", "apply"},
			Usage:   "Apply snapshot to target instance",
			Action:  cmd.ApplySchema,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:      "input",
					Aliases:   []string{"i"},
					Usage:     "Snapshot file as input",
					EnvVars:   []string{"DRCTS_SCHEMA_FILE"},
					TakesFile: true,
				},
				&cli.BoolFlag{
					Name:        "force",
					Aliases:     []string{"f"},
					DefaultText: "false",
					Usage:       "Ignore version and db differences between instances",
					EnvVars:     []string{"DRCTS_FORCE"},
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"v"},
					DefaultText: "false",
					Usage:       "Verbose",
				},
				&cli.BoolFlag{
					Name:        "yes",
					Aliases:     []string{"y"},
					DefaultText: "false",
					Usage:       "Don't ask for confirm",
				},
			},
		},
	},
}
