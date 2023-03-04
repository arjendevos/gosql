package functions

import (
	"os"
	"os/exec"
)

func (c *GoSQLConfig) SetupApi(models []*Model) error {
	if !c.SetupProject {
		return nil
	}
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

	if _, err := os.Stat(projectDir + "/database/client.go"); os.IsNotExist(err) {
		if err := populateTemplate("templates/setup/database.gotpl", projectDir+"/database/client.go", SetupMainTemplateData{PackageName: "database", Imports: imports}); err != nil {
			return err
		}
	}

	if _, err := os.Stat(projectDir + "/migrate/migrate.go"); os.IsNotExist(err) {
		if err := populateTemplate("templates/setup/migrate.gotpl", projectDir+"/database/migrate.go", SetupMainTemplateData{PackageName: "database", Imports: imports, MigrationPath: c.MigrationDir}); err != nil {
			return err
		}
	}

	imports = append(imports, moduleName+"/database")

	hasAlreadyFullSetup := false
	if _, err := os.Stat(projectDir + "/database/main.go"); os.IsNotExist(err) {
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
		if _, err := os.Stat(projectDir + "/auth/calls.go"); os.IsNotExist(err) {
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

		if err := populateTemplate("templates/setup/main.gotpl", projectDir+"/main.go", SetupMainTemplateData{PackageName: "main", Imports: imports, FullSetup: true, HasExtraMiddleWare: hasAuthUser && hasAuthOrganization && hasAuthOrganizationUser}); err != nil {
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
