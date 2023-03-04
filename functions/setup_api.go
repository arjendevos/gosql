package functions

import (
	"os"
	"os/exec"
)

func (c *GoSQLConfig) SetupApi(models []*Model) error {
	projectDir, _ := os.Getwd()

	modFile, err := parseModFile()
	if err != nil {
		return err
	}

	moduleName := modFile.Module.Mod.Path

	var imports []string

	var hasAuthUser bool
	var hasAuthOrganization bool
	var hasAuthOrganizationUser bool
	for _, model := range models {
		if model.IsAuthUser {
			hasAuthUser = true
		}
		if model.IsAuthOrganization {
			hasAuthOrganization = true
		}
		if model.IsAuthOrganizationUser {
			hasAuthOrganizationUser = true
		}
	}

	_, err2 := os.Stat(projectDir + "/database/client.go")
	if err2 != nil {
		if err := populateTemplate("templates/setup/database.gotpl", projectDir+"/database/client.go", SetupMainTemplateData{PackageName: "database", Imports: imports}); err != nil {
			return err
		}
	}

	_, err3 := os.Stat(projectDir + "/database/client.go")
	if err3 != nil {
		if err := populateTemplate("templates/setup/database.gotpl", projectDir+"/database/client.go", SetupMainTemplateData{PackageName: "database", Imports: imports}); err != nil {
			return err
		}
	}

	_, err4 := os.Stat(projectDir + "/migrate/migrate.go")
	if err4 != nil {
		if err := populateTemplate("templates/setup/migrate.gotpl", projectDir+"/database/migrate.go", SetupMainTemplateData{PackageName: "database", Imports: imports}); err != nil {
			return err
		}
	}

	imports = append(imports, moduleName+"/database")

	hasAlreadyFullSetup := false
	_, err5 := os.Stat(projectDir + "/database/main.go")
	if err5 != nil {
		if err := populateTemplate("templates/setup/main.gotpl", projectDir+"/main.go", SetupMainTemplateData{PackageName: "main", Imports: imports, FullSetup: false}); err != nil {
			return err
		}
	} else {
		hasAlreadyFullSetup = true
	}

	if !hasAlreadyFullSetup {
		if err := migrate(); err != nil {
			return err
		}
	}

	if hasAuthUser && hasAuthOrganization && hasAuthOrganizationUser {
		_, err := os.Stat(projectDir + "/auth/calls.go")
		if err != nil {
			if err := os.MkdirAll("auth", os.ModePerm); err != nil {
				return err
			}
			var imports2 []string
			imports2 = addImport(imports2, moduleName+"/"+c.ControllerOutputDir)
			if err := populateTemplate("templates/setup/auth_calls.gotpl", projectDir+"/auth/calls.go", SetupMainTemplateData{PackageName: "auth", Imports: imports2}); err != nil {
				return err
			}
		}

		imports = append(imports, moduleName+"/auth")
	}

	if !hasAlreadyFullSetup {
		imports = addImport(imports, moduleName+"/"+c.ControllerOutputDir)

		if err := populateTemplate("templates/setup/main.gotpl", projectDir+"/main.go", SetupMainTemplateData{PackageName: "main", Imports: imports, FullSetup: true}); err != nil {
			return err
		}
	}

	return nil
}

func migrate() error {
	goModTidy := exec.Command("/opt/homebrew/bin/go", []string{"mod", "tidy"}...)
	err := goModTidy.Run()
	if err != nil {
		return err
	}

	cmd := exec.Command("/opt/homebrew/bin/go", []string{"run", "main.go", "--migrate"}...)
	return cmd.Run()
}
