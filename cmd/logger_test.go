package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"os"
	"testing"
)

func Test_PrintWarn(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintWarn("This is a test!")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	wantMsg := fmt.Sprintln(aurora.Red("✖ This is a test!"))
	if string(out) != wantMsg {
		t.Errorf("%#v, wanted %#v", string(out), wantMsg)
	}
}

func Test_PrintSuccess(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintSuccess("This is a test!")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	wantMsg := fmt.Sprintln(aurora.Green("✔ This is a test!"))
	if string(out) != wantMsg {
		t.Errorf("%#v, wanted %#v", string(out), wantMsg)
	}
}

func Test_PrintInfo(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintInfo("This is a test!")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	wantMsg := fmt.Sprintln(aurora.Gray(14, "✔ This is a test!"))
	if string(out) != wantMsg {
		t.Errorf("%#v, wanted %#v", string(out), wantMsg)
	}
}
