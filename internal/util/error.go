package util

import (
	"database/sql"
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckIfNoRowsFound(err error, message string) {
	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		os.Exit(1)
	} else {
		log.Fatal(err)
		os.Exit(1)
	}
}
