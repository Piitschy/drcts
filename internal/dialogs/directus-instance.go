package dialogs

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/Piitschy/drcts/internal/directus"
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
			if err = d.TestConnection(); err != nil {
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
	email := cCtx.String(prefix + "-email")
	password := cCtx.String(prefix + "-password")
	method := ""

	if token != "" {
		method = "Token"
	} else if email != "" || password != "" {
		method = "Email/Password"
	}

	if email == "" && password == "" && token == "" {
		methodPrompt := promptui.Select{
			Label: "Choose authentication method",
			Items: []string{"Token", "Email/Password"},
		}
		var err error
		_, method, err = methodPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return nil, err
		}
	}

	if token == "" && method == "Token" {
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

	if email == "" && method == "Email/Password" {
		prompt := promptui.Prompt{
			Label:    "Email of the " + prefix + " Directus instance",
			Validate: validateEmail,
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return nil, err
		}
		email = result
	}

	if password == "" && method == "Email/Password" {
		prompt := promptui.Prompt{
			Label:    "Password of the " + prefix + " Directus instance",
			Mask:     '*',
			Validate: validatePassword,
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return nil, err
		}
		password = result
	}

	d, err := directus.NewDirectus(url, token)
	if err != nil {
		fmt.Printf("Failed to create Directus instance %v\n", err)
		return nil, err
	}
	if method == "Email/Password" {
		err = d.Login(email, password)
		if err != nil {
			fmt.Printf("Login failed %v\n", err)
			return nil, err
		}
	}
	err = d.TestConnection()
	if err != nil {
		fmt.Printf("Connection failed %v\n", err)
		return nil, err
	}
	return d, nil
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

func validateEmail(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("Email cannot be empty")
	}

	if !strings.Contains(input, "@") {
		return fmt.Errorf("Email must contain an @")
	}
	_, err := mail.ParseAddress(input)
	if err != nil {
		return fmt.Errorf("Invalid email address")
	}
	return nil
}

func validatePassword(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("Password cannot be empty")
	}
	return nil
}
