package dialogs

import (
	"fmt"
	"strings"

	"github.com/Piitschy/drctsdm/internal/directus"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

// DirectusInstance prompts the user for the Directus instance URL and token
func DirectusInstance(cCtx *cli.Context, prefix string) (*directus.Directus, error) {
	url := cCtx.String(prefix + "-url")
	for {
		if url == "" {
			prompt := promptui.Prompt{
				Label:     "URL of the " + prefix + " Directus instance",
				Validate:  validateURL,
				Default:   "http://localhost:8055",
				AllowEdit: true,
			}
			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return nil, err
			}
			d, err := directus.NewDirectus(result, "")
			if err != nil {
				fmt.Printf("Invalid URL %v\n", err)
				continue
			}
			err = d.TestConnection()
			if err != nil {
				fmt.Printf("Connection failed %v\n", err)
				continue
			}
			url = result
		}
		if url != "" {
			break
		}
	}

	token := cCtx.String(prefix + "-token")
	if token == "" {
		prompt := promptui.Prompt{
			Label:    "Token of the " + prefix + " Directus instance",
			Mask:     '*',
			Validate: validateToken,
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return nil, err
		}
		token = result
	}

	return directus.NewDirectus(url, token)
}

func validateURL(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("URL cannot be empty")
	}
	if !strings.HasPrefix(input, "http") {
		return fmt.Errorf("URL must start with http(s)://")
	}
	return nil
}

func validateToken(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("Token cannot be empty")
	}
	return nil
}
