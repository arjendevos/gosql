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
	SetupProject        bool
	TypesOutputDir      string
}

func Convert(c *GoSQLConfig) error {
	schemeDirEndsWithSlash := strings.HasSuffix(c.SchemeDir, "/")
	if !schemeDirEndsWithSlash {
		c.SchemeDir += "/"
	}

	var newFiles []string
	files, err := filepath.Glob(c.SchemeDir + "**.gosql")
	if err != nil {
		return err
	}

	var sqlType string
	var models []*Model

	err = os.RemoveAll(c.MigrationDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(c.MigrationDir); os.IsNotExist(err) {
		fmt.Println("[1] -- No migration directory found")
		if err := os.MkdirAll(c.MigrationDir, os.ModePerm); err != nil {
			return err
		}
		for _, filePath := range files {
			fileName := strings.TrimSuffix(strings.TrimPrefix(filePath, c.SchemeDir), ".gosql")
			if _, err := os.Stat(c.MigrationDir + "/" + fileName + ".up.sql"); os.IsNotExist(err) {
				newFiles = append(newFiles, filePath)
			}
		}
	} else {
		fmt.Println("[1] -- Migration directory found")

		for _, filePath := range files {
			fileName := strings.TrimSuffix(strings.TrimPrefix(filePath, c.SchemeDir), ".gosql")
			if _, err := os.Stat(c.MigrationDir + "/" + fileName + ".up.sql"); os.IsNotExist(err) {
				newFiles = append(newFiles, filePath)
			}
		}
	}

	// parse all files to get all models
	var cache = make(map[string][]*Model)
	for _, filePath := range files {
		t, mdls := ParseGoSQLFile(filePath, sqlType)
		if sqlType == "" {
			sqlType = t
		}
		models = append(models, mdls...)
		cache[filePath] = mdls
	}

	if len(newFiles) > 0 {
		err = os.WriteFile(c.MigrationDir+"/1_migration.up.sql", []byte("DROP SCHEMA IF EXISTS public CASCADE; \nCREATE SCHEMA public; CREATE table schema_migrations (version varchar(255) NOT NULL,\ndirty boolean NOT NULL DEFAULT false,\nPRIMARY KEY (version));"), 0644)
		if err != nil {
			return err
		}

		fmt.Println("[2] -- Found " + fmt.Sprint(len(newFiles)) + " files to convert")

		for _, filePath := range newFiles {
			s := strings.Split(filePath, "/")
			fileName := s[len(s)-1]

			err = c.ConvertToSql(fileName, sqlType, cache[filePath], models)
			if err != nil {
				return fmt.Errorf("error converting to sql: %v", err)
			}
		}

		// fmt.Println("[3] -- Migrating converted sql files")
		// if err := migrate(); err != nil {
		// 	return fmt.Errorf("error migrating: ", err)
		// }

	} else {
		fmt.Println("[2] -- No new files found, skipping sql conversion")
		// fmt.Println("[3] -- Nothing to migrate, skipping migration")
	}

	fmt.Println("[3] -- Setting up api & migrating (if needed)")
	err = c.InitialSetup(models)
	if err != nil {
		return fmt.Errorf("error setting up api: %v", err)
	}

	if err := migrate(); err != nil {
		return err
	}

	fmt.Println("[4] -- Converting models")
	if err := c.ConvertApiModels(models); err != nil {
		return fmt.Errorf("error converting models: %v", err)
	}

	fmt.Println("[5] -- Converting controllers")
	if err := c.ConvertApiControllers(models); err != nil {
		return fmt.Errorf("error converting controllers: %v", err)
	}

	fmt.Println("[6] -- Converting typescript types")
	if err := c.ConvertTypes(models); err != nil {
		return fmt.Errorf("error converting typescript types: %v", err)
	}

	fmt.Println("[7] -- Downloading dependencies")
	goModTidy := exec.Command("/opt/homebrew/bin/go", []string{"mod", "tidy"}...)
	if err := goModTidy.Run(); err != nil {
		return fmt.Errorf("error running go mod tidy: %v", err)
	}

	fmt.Println("[8] -- Converted all files successfully")
	return nil
}
