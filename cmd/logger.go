package cmd

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
)

func PrintWarn(s string) {
	fmt.Println(Sprintf(Red("✖ %s"), s))
}

func PrintSuccess(s string) {
	fmt.Println(Sprintf(Green("✔ %s"), s))
}

func PrintInfo(s string) {
	fmt.Println(Sprintf(Gray(14, "➜ %s"), s))
}
