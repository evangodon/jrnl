package util

import (
	"testing"
	"time"
)

func TestFormatToLocalTime(t *testing.T) {
	utcDate := time.Date(2014, time.November, 28, 1, 1, 1, 1, time.UTC)

	got := FormatToLocalTime(utcDate, "2006-01-02")
	want := "2014-11-27"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}
