package gosql

import (
	"os"

	"github.com/arjendevos/gosql/functions"
	"github.com/gin-gonic/gin"
)

type GoSQLConfig struct {
	SchemeDir           string
	MigrationDir        string
	ModelOutputDir      string
	ControllerOutputDir string
	AuthMiddelware      func(g *gin.Context)
}

func Convert(c *GoSQLConfig) {
	dir, _ := os.Getwd()

	functions.Convert(&functions.GoSQLConfig{
		SchemeDir:           dir + "/" + c.SchemeDir,
		MigrationDir:        dir + "/" + c.MigrationDir,
		ModelOutputDir:      c.ModelOutputDir,
		ControllerOutputDir: c.ControllerOutputDir,
		AuthMiddelware:      c.AuthMiddelware,
	})
}
