package gosql

import (
	"github.com/arjendevos/gosql/functions"
)

type GoSQLConfig struct {
	SchemeDir           string
	MigrationDir        string
	ModelOutputDir      string
	ControllerOutputDir string
}

func Convert(c *GoSQLConfig) {
	functions.Convert(&functions.GoSQLConfig{
		SchemeDir:           c.SchemeDir,
		MigrationDir:        c.MigrationDir,
		ModelOutputDir:      c.ModelOutputDir,
		ControllerOutputDir: c.ControllerOutputDir,
	})
}
