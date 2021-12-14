package util

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func GetNewEntry(content string) string {
	file, err := os.CreateTemp("", "today-*.md")
	CheckError(err)

	contentToWrite := content
	if contentToWrite == "" {
		todayDate := time.Now().Format("Monday, January 2 2006")
		contentToWrite = "# " + todayDate + "\n\n"
	}

	_, err = file.WriteString(contentToWrite)
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
