package lib

import (
	"io/ioutil"
	"os"
	"os/exec"
)

const DefaultEditor = "vim"

// OpenFileInEditor opens filename in a text editor and obeys the settings
// via ENV variable `EDITOR`. Defaults to `vim` as the editor if ENV variable
// is not set.
func OpenFileInEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	// Get the full executable path for the editor.
	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetInputFromEditor opens a temporary file in a text editor and returns
// the written bytes on success or an error on failure. Temp file will be
// deleted automatically.
func GetInputFromEditor(data []byte) ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return []byte{}, err
	}

	// Write data if we have it
	if len(data) > 0 {
		file.Write(data)
	}

	filename := file.Name()

	// Defer removal of the temporary file in case any of the next steps fail.
	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = OpenFileInEditor(filename); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
