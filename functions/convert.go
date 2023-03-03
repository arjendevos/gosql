package functions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GoSQLConfig struct {
	SchemeDir           string
	MigrationDir        string
	ModelOutputDir      string
	ControllerOutputDir string
}

func Convert(c *GoSQLConfig) {
	schemeDirEndsWithSlash := strings.HasSuffix(c.SchemeDir, "/")
	if !schemeDirEndsWithSlash {
		c.SchemeDir += "/"
	}

	files, err := filepath.Glob(c.SchemeDir + "**.gosql")
	if err != nil {
		fmt.Println("ERR!", err)
		os.Exit(1)
	}

	for _, filePath := range files {
		s := strings.Split(filePath, "/")
		fileName := s[len(s)-1]
		sqlType, models := ParseGoSQLFile(filePath)
		err := c.ConvertToSql(fileName, sqlType, models)
		if err != nil {
			fmt.Println("ERR!", err)
			os.Exit(1)
		}

		err = c.ConvertApiModels(models)
		if err != nil {
			fmt.Println("ERR!", err)
			os.Exit(1)
		}

		err = c.ConvertApiControllers(models)
		if err != nil {
			fmt.Println("ERR!", err)
			os.Exit(1)
		}
	}
}
