package main

import (
	"flag"
	"fmt"

	"github.com/Ak-Army/git-prompt/color"
	"github.com/Ak-Army/git-prompt/prompt"
)

func main() {
	status := prompt.GetCurrentStatus()

	var fColored = flag.String("colored", "default", "colored library (default, zsh)")
	flag.Parse()

	var colored color.Color
	if *fColored == "zsh" {
		colored = color.Color{
			color.ZshColoredOutput{},
		}
	} else {
		colored = color.Color{
			color.DefaultColoredOutput{},
		}
	}
	fmt.Printf("%+v",status.Prompt(colored))
}
