package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/evangodon/jrnl/internal/db"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {
	dat, err := os.ReadFile("./jrnl/work-journal.txt")
	check(err)

	entries := strings.Split(string(dat), "[")

	dict := map[string]db.Journal{}

	for _, v := range entries {
		if len(v) > 25 {
			createdStr := v[0:22]
			createdStr = strings.Replace(createdStr, "[", "", 2)
			createdStr = strings.Replace(createdStr, "]", "", 2)
			createdStr = strings.Replace(createdStr, "-", " ", 2)
			content := v[23:]

			createdAt, err := time.Parse("Mon Jan 02 2006 03:04", createdStr)
			check(err)

			createdAtFormatted := createdAt.Format("2006-01-02")

			if dict[createdAtFormatted].Content != "" {
				newContent := dict[createdAt.Format("2006-01-02")].Content + content
				updatedDaily := db.Journal{
					ID:        db.CreateID(),
					CreatedAt: createdAt,
					UpdatedAt: createdAt,
					Content:   newContent,
				}
				dict[createdAtFormatted] = updatedDaily
			} else {
				dict[createdAtFormatted] = db.Journal{
					ID:        db.CreateID(),
					CreatedAt: createdAt,
					UpdatedAt: createdAt,
					Content:   content,
				}
			}
		}
	}

	dbClient := db.Connect()
	ctx := context.Background()
	for _, daily := range dict {
		println(daily.CreatedAt.String(), daily.Content, "\n")
		_, err = dbClient.NewInsert().
			Model(&daily).
			Exec(ctx)
		check(err)
	}

}
