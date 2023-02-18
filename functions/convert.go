package functions

import (
	"fmt"
	"path/filepath"
)

type GoSQLConfig struct {
	SchemeDir           string
	MigrationDir        string
	ModelOutputDir      string
	ControllerOutputDir string
}

func Convert(c *GoSQLConfig) {
	files, err := filepath.Glob(c.SchemeDir + "/*.gosql")
	if err != nil {
		panic(err)
	}

	for _, fileName := range files {
		fmt.Println(fileName)
		// sqlType, models := parseGoSQLFile(fileName)
		// err := c.ConvertToSql(fileName, sqlType, models)
		// if err != nil {
		// 	panic(err)
		// }

		// err = c.ConvertApiModels(models)
		// if err != nil {
		// 	panic(err)
		// }

		// err = c.ConvertApiControllers(models)
		// if err != nil {
		// 	panic(err)
		// }
	}
}
