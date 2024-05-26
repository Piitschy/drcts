package dialogs

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func YesNo(cCtx *cli.Context, question string, a ...any) (bool, error) {
	y := false
	if cCtx.Bool("yes") {
		y = true
		return y, nil
	}
	question = fmt.Sprintf(question, a...)
	prompt := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	y = strings.ToLower(result) == "y" || strings.ToLower(result) == "yes" || strings.ToLower(result) == "Y"

	if !y {
		fmt.Println("Aborted")
	}
	return y, nil
}
