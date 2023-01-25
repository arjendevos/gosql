package main

import (
	"path/filepath"

	"github.com/arjendevos/gosql/functions"
)

func main() {
	files, _ := filepath.Glob("*.gosql")

	for _, file := range files {
		sqlType, models := functions.ParseGoSQLFile(file)
		if sqlType == "postgresql" {
			functions.ConvertPostgreSQL(file, models)
		}
	}

	// sqlType, models := functions.ParseGoSQLFile("test")
	// if sqlType == "postgresql" {
	// 	functions.ConvertPostgreSQL("models", models)
	// }

}
