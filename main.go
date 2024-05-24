package main

import (
	"fmt"

	"github.com/Piitschy/drctsdm/internal/directus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "drctsdm",
		Usage: "Directus Data Model CLI for schema migration",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	_ = app
	// if err := app.Run(os.Args); err != nil {
	// 	log.Fatal(err)
	// }

	d, err := directus.NewDirectus("http://localhost:8055", "4PopL2EUW67BqsabSrWODHN7ombpp3rH")

	if err != nil {
		fmt.Println(err)
	}

}

func MigrateSchema(cCtx *cli.Context) error {
	dBase, err := directus.NewDirectus("http://localhost:8055", "4PopL2EUW67BqsabSrWODHN7ombpp3rH")
	dTarget, err := directus.NewDirectus("http://localhost:8055", "4PopL2EUW67BqsabSrWODHN7ombpp3rH")
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
