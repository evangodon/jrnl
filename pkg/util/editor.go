package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// GetEditor creates a temporary file and opens it with the default editor,
// and then returns the content written to that file.
func GetNewEntry(content string) string {
	file, err := os.CreateTemp("", "today-*.md")
	CheckError(err)

	_, err = file.WriteString(content)
	CheckError(err)

	OpenNoteWithEditor(file.Name())

	newContent, err := os.ReadFile(file.Name())
	CheckError(err)

	defer os.Remove(file.Name())

	return string(newContent)
}

// OpenEditor opens the default editor with the given filename.
func OpenNoteWithEditor(notePath string) {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		fmt.Println("Error: no editor found in environment")
		os.Exit(1)
	}

	e := exec.Command(editor, notePath)

	e.Stdin = os.Stdin
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr

	err := e.Run()

	if err != nil {
		log.Fatal(err)
	}
}
