package functions

import (
	"fmt"
	"os"
	"os/exec"
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

	var sqlType string
	var models []*Model

	if err := os.RemoveAll(c.MigrationDir); err != nil {
		fmt.Println("ERR!", err)
		os.Exit(2)
	}
	if err := os.MkdirAll(c.MigrationDir, os.ModePerm); err != nil {
		fmt.Println("ERR!", err)
		os.Exit(3)
	}

	for _, filePath := range files {
		s := strings.Split(filePath, "/")
		fileName := s[len(s)-1]
		t, mdls := ParseGoSQLFile(filePath, sqlType)
		if sqlType == "" {
			sqlType = t
		}

		models = append(models, mdls...)

		err = c.ConvertToSql(fileName, sqlType, mdls, models)
		if err != nil {
			fmt.Println("ERR!", err)
			os.Exit(4)
		}
	}

	err = c.SetupApi(models)
	if err != nil {
		fmt.Println("ERR!", err)
		os.Exit(5)
	}

	err = c.ConvertApiModels(models)
	if err != nil {
		fmt.Println("ERR!", err)
		os.Exit(6)
	}

	err = c.ConvertApiControllers(models)
	if err != nil {
		fmt.Println("ERR!", err)
		os.Exit(7)
	}

	err = c.ConvertTypes(models)
	if err != nil {
		fmt.Println("ERR!", err)
		os.Exit(8)
	}

	goModTidy := exec.Command("/opt/homebrew/bin/go", []string{"mod", "tidy"}...)
	err = goModTidy.Run()
	if err != nil {
		fmt.Println("ERR!", err)
		os.Exit(9)
	}
}
