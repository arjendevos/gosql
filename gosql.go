package gosql

import (
	"fmt"
	"os"

	"github.com/arjendevos/gosql/functions"
)

type GoSQLConfig struct {
	SchemeDir           string
	MigrationDir        string
	ModelOutputDir      string
	ControllerOutputDir string
	SetupProject        bool
	TypesOutputDir      string
}

func Convert(c *GoSQLConfig) {
	dir, _ := os.Getwd()

	err := functions.Convert(&functions.GoSQLConfig{
		SchemeDir:           dir + "/" + c.SchemeDir,
		MigrationDir:        dir + "/" + c.MigrationDir,
		ModelOutputDir:      c.ModelOutputDir,
		ControllerOutputDir: c.ControllerOutputDir,
		SetupProject:        c.SetupProject,
		TypesOutputDir:      c.TypesOutputDir,
	})
	if err != nil {
		fmt.Println("ERR", err)
		os.Exit(1)
	}
}
