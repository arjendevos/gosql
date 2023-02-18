package main

import (
	"github.com/arjendevos/gosql/functions"
)

func main() {
	functions.Convert(&functions.GoSQLConfig{
		SchemeDir:           "/",
		MigrationDir:        "/database/migrations",
		ModelOutputDir:      "/models",
		ControllerOutputDir: "/controllers",
	})

}
