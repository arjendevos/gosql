package main

import (
	"path/filepath"

	"github.com/arjendevos/gosql/functions"
)

func main() {
	files, _ := filepath.Glob("*.gosql")

	for _, file := range files {
		_, models := functions.ParseGoSQLFile(file)
		// err := functions.ConvertToSql(file, sqlType, models)
		// if err != nil {
		// 	panic(err)
		// }

		err := functions.ConvertToApiModels(models)
		if err != nil {
			panic(err)
		}

		err = functions.ConvertToApiControllers(models)
		if err != nil {
			panic(err)
		}
	}

	// sqlType, models := functions.ParseGoSQLFile("test")
	// if sqlType == "postgresql" {
	// 	functions.ConvertPostgreSQL("models", models)
	// }

}
