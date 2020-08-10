package cmd

import (
	"fmt"
	a "github.com/logrusorgru/aurora"
)

func PrintWarn(s string) {
	fmt.Println(a.Sprintf(a.Red("✖ %s"), s))
}

func PrintSuccess(s string) {
	fmt.Println(a.Sprintf(a.Green("✔ %s"), s))
}

func PrintInfo(s string) {
	fmt.Println(a.Sprintf(a.Gray(14, "➜ %s"), s))
}
