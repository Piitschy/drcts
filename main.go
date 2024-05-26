package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Piitschy/drctsdm/internal/dialogs"
	"github.com/Piitschy/drctsdm/internal/directus"
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
				Usage:   "Migrate schema from base to target instance",
				Action:  Migrate,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						DefaultText: "false",
						Usage:       "Force apply diff",
						EnvVars:     []string{"FORCE"},
					},
				},
			},
			{
				Name:    "save-snapshot",
				Aliases: []string{"ss", "save"},
				Usage:   "Save snapshot of the base instance to output file (-o) or stdout",
				Action:  SaveSnapshot,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "output",
						Aliases:   []string{"o"},
						Usage:     "Output file",
						EnvVars:   []string{"OUTPUT_FILE"},
						TakesFile: true,
					},
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"F"},
						Usage:   "Output format (json, yaml, xml, csv)",
						EnvVars: []string{"OUTPUT_FORMAT"},
					},
				},
			},
			{
				Name:    "save-diff",
				Aliases: []string{"sd", "diff"},
				Usage:   "Save diff of the base and target instances to output file (-o) or stdout",
				Action:  SaveDiff,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "input",
						Aliases:   []string{"i"},
						Usage:     "Give a schema file to compare with the target instance",
						EnvVars:   []string{"BASE_SCHEMA_FILE"},
						TakesFile: true,
					},
					&cli.StringFlag{
						Name:      "output",
						Aliases:   []string{"o"},
						Usage:     "Output file",
						EnvVars:   []string{"DIFF_FILE"},
						TakesFile: true,
					},
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"F"},
						Usage:   "Output format (json, yaml, xml, csv)",
						EnvVars: []string{"DIFF_FILE_FORMAT"},
					},
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						DefaultText: "false",
						Usage:       "Force apply diff",
						EnvVars:     []string{"FORCE"},
					},
				},
			},
			{
				Name:    "apply-diff",
				Aliases: []string{"ad"},
				Usage:   "Apply diff to target instance",
				Action:  ApplyDiff,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "input",
						Aliases:   []string{"i"},
						Usage:     "Input file",
						EnvVars:   []string{"DIFF_FILE"},
						TakesFile: true,
					},
				},
			},
			{
				Name:    "apply-snapshot",
				Aliases: []string{"as", "apply"},
				Usage:   "Apply snapshot to target instance",
				Action:  ApplySnapshot,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "input",
						Aliases:   []string{"i"},
						Usage:     "Snapshot file as input",
						EnvVars:   []string{"TARGET_SCHEMA_FILE"},
						TakesFile: true,
					},
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						DefaultText: "false",
						Usage:       "Force apply diff",
						EnvVars:     []string{"FORCE"},
					},
				},
			},
		},
	}

	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func Migrate(cCtx *cli.Context) error {
	dBase, err := dialogs.DirectusInstance(cCtx, "base")
	dTarget, err := dialogs.DirectusInstance(cCtx, "target")
	if err != nil {
		return err
	}

	force := cCtx.Bool("force")

	s, err := dBase.GetSnapshot()
	diff, err := dTarget.GetDiff(s, force)

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

func SaveSnapshot(cCtx *cli.Context) error {
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

	b, err := dBase.GetRawSnapshot(fileExt[1:])

	if err != nil {
		return err
	}

	if cCtx.String("output") == "" {
		os.Stdout.Write(b)
		return nil
	}
	fmt.Println("Writing snapshot to " + cCtx.String("output"))
	os.WriteFile(cCtx.String("output"), b, 0644)

	return nil
}

func ApplySnapshot(cCtx *cli.Context) error {
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

	if err != nil {
		return err
	}

	return dTarget.ApplyDiff(diff)
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
	fmt.Println("Writing diff to " + cCtx.String("output"))
	os.WriteFile(cCtx.String("output"), diffBytes, 0644)

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

	err = dTarget.ApplyDiff(&diff)
	if err != nil {
		return err
	}

	return nil
}
