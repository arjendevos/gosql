package main

import (
	"github.com/arjendevos/gosql/functions"
)

func main() {
	functions.Convert(&functions.GoSQLConfig{
		SchemeDir:           "test",
		MigrationDir:        "database/migrations",
		ModelOutputDir:      "models",
		ControllerOutputDir: "controllers",
	})

}
