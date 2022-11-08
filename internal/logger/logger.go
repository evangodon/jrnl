package logger

import (
	"io"
	"os"
)

type Logger struct {
	out io.Writer
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{out}
}

func (l Logger) Print(msg string) {
	l.out.Write([]byte(msg))
}

func (l Logger) PrintInfo(msg string) {
	log := " ℹ️ " + msg + "\n"
	l.out.Write([]byte(log))
}

func (l Logger) PrintSuccess(msg string) {
	log := " ✅ " + msg + "\n"
	l.out.Write([]byte(log))
}

func (l Logger) PrintError(msg string) {
	log := " ❗ " + msg + "\n"
	l.out.Write([]byte(log))
}

func (l Logger) PrintFatal(msg string) {
	log := " ❌ " + msg + "\n"
	l.out.Write([]byte(log))

	l.out.Write([]byte("Exiting...\n"))
	os.Exit(1)
}
