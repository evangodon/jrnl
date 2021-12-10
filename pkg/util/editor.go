package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func OpenEditor() string {
	content := ""
	file, err := os.CreateTemp("", "today-*.md")

	CheckError(err)

	OpenNoteWithEditor(file.Name())

	file.Seek(0, 0)
	s := bufio.NewScanner(file)

	for s.Scan() {
		content += s.Text()
	}

	defer os.Remove(file.Name())

	return content
}

// OpenEditor opens the default editor with the given filename.
func OpenNoteWithEditor(notePath string) {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		fmt.Println("Error: no editor found in environment")
		os.Exit(1)
	}

	editorCmd := exec.Command(editor, notePath)

	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	err := editorCmd.Start()

	if err != nil {
		panic(err)
	}

	err = editorCmd.Wait()

	if err != nil {
		panic(err)
	}
}
