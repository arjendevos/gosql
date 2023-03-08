package functions

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"unicode"
)

type StructWithRelations struct {
	StructName string
	Relations  []string
}

func (c *GoSQLConfig) ConvertApiModels(models []*Model) error {
	config, err := parseSqlBoilerConfig()
	if err != nil {
		return err
	}

	modFile, err := parseModFile()
	if err != nil {
		return err
	}

	moduleName := modFile.Module.Mod.Path

	pkgName := config.PkgName
	dir, _ := os.Getwd()
	outputDir := dir + "/" + c.ModelOutputDir + "/am"

	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	if err := populateTemplate("templates/model/model_helpers.gotpl", outputDir+"/helpers.go", GeneralTemplateData{
		PackageName: "am",
	}); err != nil {
		return err
	}

	for _, m := range models {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fmt.Sprintf("%s/%s.go", config.Output, m.SnakeName), nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("make sure to migrate and run sqlboiler first")
		}

		var imports []string
		imports = addImport(imports, moduleName+"/"+config.Output) // sql boiler models

		relations := getRelations(f, m, pkgName, models)

		var columnsWithRelationsAsIDs []*Column
		for _, cl := range m.Columns {
			if !cl.IsRelation {
				columnsWithRelationsAsIDs = append(columnsWithRelationsAsIDs, cl)
				// if cl.Type.GoTypeName == "time.Time" {
				// 	imports = addImport(imports, "time")
				// }
			}

			if cl.IsRelation {
				columnsWithRelationsAsIDs = append(columnsWithRelationsAsIDs, &Column{
					SnakeName:  cl.SnakeName + "_id",
					CamelName:  cl.CamelName + "ID",
					Type:       &Type{Name: "int", GoTypeName: "int", IsNullable: cl.Type.IsNullable, EmptyValue: cl.Type.EmptyValue},
					Attributes: cl.Attributes,
					IsRelation: true,
					Expose:     true,
				})
			}
		}

		if err := populateTemplate("templates/model/model.gotpl", outputDir+"/"+m.SnakeName+".go", ModelTemplateData{
			PackageName: "am",
			CamelName:   m.CamelName,
			Imports:     imports,
			Relations:   relations,
			Columns:     columnsWithRelationsAsIDs,
		}); err != nil {
			return err
		}
	}

	return nil

}

func (c *GoSQLConfig) ConvertApiControllers(models []*Model) error {
	dir, _ := os.Getwd()

	outputDir := dir + "/" + c.ControllerOutputDir
	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	config, err := parseSqlBoilerConfig()
	if err != nil {
		return err
	}

	modFile, err := parseModFile()
	if err != nil {
		return err
	}

	moduleName := modFile.Module.Mod.Path
	pkgName := config.PkgName
	var imports []string

	imports = addImport(imports, moduleName+"/"+config.Output)
	imports = addImport(imports, moduleName+"/"+c.ModelOutputDir+"/am")

	if err := populateTemplate("templates/api/columns.gotpl", outputDir+"/generated_columns.go", GeneralTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: models}); err != nil {
		return err
	}

	if err := populateTemplate("templates/api/orders.gotpl", outputDir+"/generated_orders.go", GeneralTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: models}); err != nil {
		return err
	}

	// Build the controllers
	var modelWithRelations []*ModelWithRelations
	var modelWithRelationsAsIdsInColumns []*Model
	var createAndUpdateData []*CreateAndUpdateDataModel
	var jwtFields []*JWTField
	var authQueryFields []*JWTField
	var hasNullableFields bool

	var authUser *Model
	var authOrganization *Model
	var authOrganizationUser *Model

	var modelImports []string
	modelImports = addImport(modelImports, moduleName+"/"+config.Output)

	for _, m := range models {
		if m.Hide {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fmt.Sprintf("%s/%s.go", config.Output, m.SnakeName), nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("make sure you run sqlboiler first")
		}
		relations := getRelations(f, m, pkgName, models)
		modelWithRelations = append(modelWithRelations, &ModelWithRelations{Model: m, Relations: relations})

		createColumns, updateColumns, mImports := getCreateAndUpdateColumns(m)

		if err := populateTemplate("templates/api/controller.gotpl", outputDir+"/generated_"+m.SnakeName+"_controller.go", ControllerTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), CamelName: m.CamelName, Imports: imports, CreateColumns: createColumns, UpdateColumns: updateColumns, Model: m}); err != nil {
			return err
		}

		modelImports = append(modelImports, mImports...)
		createAndUpdateData = append(createAndUpdateData, &CreateAndUpdateDataModel{
			SnakeName:     m.SnakeName,
			CamelName:     m.CamelName,
			CreateColumns: createColumns,
			UpdateColumns: updateColumns,
			Model:         m,
		})

		var columnsWithRelationsAsIDs []*Column
		for _, c := range m.Columns {
			if !c.IsRelation {
				columnsWithRelationsAsIDs = append(columnsWithRelationsAsIDs, c)
			}

			if c.IsRelation {
				columnsWithRelationsAsIDs = append(columnsWithRelationsAsIDs, &Column{
					SnakeName:  c.SnakeName + "_id",
					CamelName:  c.CamelName + "ID",
					Type:       &Type{Name: "int", GoTypeName: "int", IsNullable: c.Type.IsNullable, EmptyValue: c.Type.EmptyValue},
					Attributes: c.Attributes,
					IsRelation: true,
					Expose:     true,
				})
			}

			if c.Type.IsNullable {
				hasNullableFields = true
			}
		}

		modelWithRelationsAsIdsInColumns = append(modelWithRelationsAsIdsInColumns, &Model{
			SnakeName:              m.SnakeName,
			CamelName:              m.CamelName,
			Columns:                columnsWithRelationsAsIDs,
			IsAuthRequired:         m.IsAuthRequired,
			IsAuthUser:             m.IsAuthUser,
			IsAuthOrganization:     m.IsAuthOrganization,
			ProtectedRoutes:        m.ProtectedRoutes,
			IsAuthOrganizationUser: m.IsAuthOrganizationUser,
		})

		if m.IsAuthUser {
			authUser = m

			for _, c := range m.Columns {
				for _, a := range c.Attributes {
					if a.Name == "unique" && !c.IsRelation {
						jwtFields = append(jwtFields, &JWTField{
							CamelName:                   m.CamelName + c.CamelName,
							SnakeName:                   m.SnakeName + "_" + c.SnakeName,
							GoType:                      c.Type.GoTypeName,
							NormalName:                  c.CamelName,
							TableCamelName:              m.CamelName,
							TableSnakeName:              m.SnakeName,
							IsFromUserTable:             true,
							IsFromOrganizationTable:     false,
							IsFromOrganizationUserTable: false,
						})
					}
				}

				if c.SnakeName == "id" {
					authQueryFields = append(authQueryFields, &JWTField{
						CamelName:                   m.CamelName + c.CamelName,
						SnakeName:                   m.SnakeName + "_" + c.SnakeName,
						GoType:                      c.Type.GoTypeName,
						NormalName:                  c.CamelName,
						TableCamelName:              m.CamelName,
						TableSnakeName:              m.SnakeName,
						IsFromUserTable:             true,
						IsFromOrganizationTable:     false,
						IsFromOrganizationUserTable: false,
					})
				}
			}
		}

		if m.IsAuthOrganization {
			authOrganization = m

			for _, c := range m.Columns {
				for _, a := range c.Attributes {
					if a.Name == "unique" && !c.IsRelation {
						jwtFields = append(jwtFields, &JWTField{
							CamelName:                   m.CamelName + c.CamelName,
							SnakeName:                   m.SnakeName + "_" + c.SnakeName,
							GoType:                      c.Type.GoTypeName,
							NormalName:                  c.CamelName,
							TableCamelName:              m.CamelName,
							TableSnakeName:              m.SnakeName,
							IsFromUserTable:             false,
							IsFromOrganizationTable:     true,
							IsFromOrganizationUserTable: false,
						})
					}
				}

				if c.SnakeName == "id" {
					authQueryFields = append(authQueryFields, &JWTField{
						CamelName:                   m.CamelName + c.CamelName,
						SnakeName:                   m.SnakeName + "_" + c.SnakeName,
						GoType:                      c.Type.GoTypeName,
						NormalName:                  c.CamelName,
						TableCamelName:              m.CamelName,
						TableSnakeName:              m.SnakeName,
						IsFromUserTable:             false,
						IsFromOrganizationTable:     true,
						IsFromOrganizationUserTable: false,
					})
				}
			}
		}

		if m.IsAuthOrganizationUser {
			authOrganizationUser = m

			for _, c := range m.Columns {
				if c.SnakeName == "role" {
					field := &JWTField{
						CamelName:                   "Role",
						SnakeName:                   "role",
						GoType:                      c.Type.GoTypeName,
						NormalName:                  "Role",
						IsFromUserTable:             false,
						IsFromOrganizationTable:     false,
						IsFromOrganizationUserTable: true,
					}
					jwtFields = append(jwtFields, field)
				}
			}
		}
	}

	createColumns, _, mImports := getCreateAndUpdateColumns(authUser)
	var createColumnsOrganization []*Column
	var createColumnsOrganizationUser []*Column
	var OrganizationCamelName string
	var OrganizationUserCamelName string

	if authUser != nil && authUser.Oauth2 != nil {
		if err := populateTemplate("templates/api/oauth2_controller.gotpl", outputDir+"/generated_oauth2_controller.go", Oauth2TemplateData{
			PackageName:           strings.ReplaceAll(c.ControllerOutputDir, "/", "_"),
			UserTable:             authUser,
			OrganizationTable:     authOrganization,
			OrganizationUserTable: authOrganizationUser,
			Imports:               modelImports,
			JWTFields:             jwtFields,
		}); err != nil {
			return err
		}
	}

	if authOrganization != nil && authOrganizationUser != nil {
		OrganizationCamelName = authOrganization.CamelName
		OrganizationUserCamelName = authOrganizationUser.CamelName

		createColumnsOrganizationModels, _, mImportsOrganization := getCreateAndUpdateColumns(authOrganization)
		createColumnsOrganizationUserModels, _, mImportsOrganizationUser := getCreateAndUpdateColumns(authOrganizationUser)

		createColumnsOrganization = createColumnsOrganizationModels
		createColumnsOrganizationUser = createColumnsOrganizationUserModels

		modelImports = append(modelImports, mImportsOrganization...)
		modelImports = append(modelImports, mImportsOrganizationUser...)
	}

	modelImports = append(modelImports, mImports...)

	if err := populateTemplate("templates/api/auth_controller.gotpl", outputDir+"/generated_auth_controller.go", AuthTemplateData{
		PackageName:                   strings.ReplaceAll(c.ControllerOutputDir, "/", "_"),
		CamelName:                     authUser.CamelName,
		Imports:                       imports,
		CreateColumns:                 createColumns,
		JWTFields:                     jwtFields,
		HasOrganization:               authOrganization != nil,
		HasOrganizationUser:           authOrganizationUser != nil,
		OrganizationCamelName:         OrganizationCamelName,
		OrganizationUserCamelName:     OrganizationUserCamelName,
		OrganizationCreateColumns:     createColumnsOrganization,
		OrganizationUserCreateColumns: createColumnsOrganizationUser,
	}); err != nil {
		return err
	}

	if err := populateTemplate("templates/api/filters.gotpl", outputDir+"/generated_filters.go", GeneralTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: models}); err != nil {
		return err
	}

	if err := populateTemplate("templates/api/queries.gotpl", outputDir+"/generated_queries.go", QueryTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: models, AuthFields: authQueryFields, HasMultipleAuthFields: len(authQueryFields) > 1}); err != nil {
		return err
	}

	if err := populateTemplate("templates/api/client.gotpl", outputDir+"/generated_client.go", QueryTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: models, AuthFields: authQueryFields}); err != nil {
		return err
	}

	if err := populateTemplate("templates/api/routes.gotpl", outputDir+"/generated_routes.go", RoutesTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: models, AuthFields: authQueryFields, HasOrganization: authOrganization != nil}); err != nil {
		return err
	}

	if err := populateTemplate("./templates/api/helpers.gotpl", outputDir+"/generated_helpers.go", HelpersTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), JWTFields: jwtFields, HasAuth: authUser != nil}); err != nil {
		return err
	}

	if authUser != nil {
		if err := populateTemplate("./templates/api/middleware.gotpl", outputDir+"/generated_middleware.go", HelpersTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), JWTFields: jwtFields, HasAuth: authUser != nil}); err != nil {
			return err
		}
	}

	if err := populateTemplate("templates/api/relations.gotpl", outputDir+"/generated_relations.go", SelectTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: modelWithRelations}); err != nil {
		return err
	}

	if err := populateTemplate("templates/api/select.gotpl", outputDir+"/generated_select.go", SelectTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: modelWithRelations}); err != nil {
		return err
	}

	bodieImports := modelImports
	if hasNullableFields {
		bodieImports = addImport(bodieImports, "github.com/volatiletech/null/v8")
	}
	if err := populateTemplate("templates/api/bodies.gotpl", outputDir+"/generated_bodies.go", BodyTemplateData{PackageName: strings.ReplaceAll(c.ControllerOutputDir, "/", "_"), Controllers: createAndUpdateData, Imports: bodieImports, AuthFields: authQueryFields}); err != nil {
		return err
	}

	return nil
}

func getCreateAndUpdateColumns(m *Model) ([]*Column, []*Column, []string) {
	var createColumns []*Column
	var updateColumns []*Column
	var imports []string

	for _, c := range m.Columns {
		if !c.IsRelation && c.SnakeName != "id" && c.SnakeName != "created_at" && c.SnakeName != "updated_at" {
			createColumns = append(createColumns, c)
			updateColumns = append(updateColumns, c)

			if c.Type.GoTypeName == "time.Time" {
				imports = addImport(imports, "time")
			}
		}

		if c.IsRelation {
			c := &Column{
				SnakeName:  c.SnakeName + "_id",
				CamelName:  c.CamelName + "ID",
				Type:       &Type{GoTypeName: "int", Name: "int", IsNullable: c.Type.IsNullable, HasDifferentCharLength: false, CharLength: 0, EmptyValue: "0"},
				Attributes: []*Attribute{},
				IsRelation: true,
				Expose:     true,
			}
			createColumns = append(createColumns, c)
			updateColumns = append(updateColumns, c)
		}
	}

	return createColumns, updateColumns, imports
}

func addImport(imports []string, importName string) []string {
	found := false
	for _, i := range imports {
		if i == importName {
			found = true
		}
	}

	if !found {
		imports = append(imports, importName)
	}

	return imports
}

func firstToLower(s string) string {
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func firstToUpper(s string) string {
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

func fieldTypeToString(prefix string, fieldType ast.Expr) (string, error) {
	switch ft := fieldType.(type) {
	case *ast.Ident:
		return prefix + ft.Name, nil
	case *ast.SelectorExpr:
		return prefix + fmt.Sprintf("%s.%s", ft.X.(*ast.Ident).Name, ft.Sel.Name), nil
	case *ast.StarExpr:
		ftStr, err := fieldTypeToString(prefix, ft.X)
		if err != nil {
			return "", err
		}
		return "*" + ftStr, nil
	default:
		return "", fmt.Errorf("unsupported field type: %T", fieldType)
	}
}

func populateTemplate(file, output string, data interface{}) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Could not get package directory path")
	}
	packageDir := path.Join(filepath.Dir(filename), "..")
	content, err := ioutil.ReadFile(packageDir + "/" + file)
	if err != nil {
		return err
	}

	template, err := parseTemplate(&TemplateConfig{Template: string(content), Data: data}, strings.HasSuffix(output, ".go"))
	if err != nil {
		return err
	}

	hF, err := os.Create(output)
	if err != nil {
		return err
	}
	defer hF.Close()
	_, err = hF.WriteString(template)
	if err != nil {
		return err
	}

	return nil
}

func getIdentifierFromAuthFields(authFields []*JWTField) *JWTField {
	for _, f := range authFields {
		if strings.EqualFold(f.NormalName, "id") {
			return f
		}
	}

	return nil
}

type TemplateConfig struct {
	Template string
	Data     interface{}
}

func parseTemplate(c *TemplateConfig, shouldFormat bool) (string, error) {
	tpl, err := template.New("").Funcs(template.FuncMap{
		"toSnake":      camelToSnake,
		"pluralize":    pluralize,
		"singularize":  singularize,
		"firstToLower": firstToLower,
		"firstToUpper": firstToUpper,
		"toLower":      strings.ToLower,
		"isFalse": func(a bool) bool {
			return !a
		},
		"everyRouteIsProtected": everyRouteIsProtected,
		"isProtected": func(r []string, rn string) bool {
			return stringArrayContains(r, rn)
		},
		"hasAuthFields": func(a []*JWTField) bool {
			return len(a) > 0
		},
		"hasOauth2": func(a []*Model) bool {
			for _, m := range a {
				if m.Oauth2 != nil {
					return true
				}
			}

			return false
		},
		"isNotNil": func(a *JWTField) bool {
			return a != nil
		},
		"isNotNilTable": func(a *Model) bool {
			return a != nil
		},
		"neq": func(a string, b string) bool {
			return !strings.EqualFold(a, b)
		},
		"contains": func(a string, b string) bool {
			return strings.Contains(a, b)
		},
		"eq": strings.EqualFold,
		"getValidate": func(c *Column) string {
			var v []string
			var canBeEmpty bool

			for _, attr := range c.Attributes {
				if attr.Name == "regexp" && attr.HasValue {
					v = append(v, "regexp="+strings.ReplaceAll(attr.Value, "'", ""))
				}

				if attr.Name == "default" {
					canBeEmpty = true
				}
			}

			if !c.Type.IsNullable && !canBeEmpty {
				v = append(v, "nonzero")
			}

			if len(v) > 0 {
				return fmt.Sprintf(" validate:\"%s\"", strings.Join(v, ","))
			}

			return ""
		},
		"isAuthFieldInModel": isAuthFieldInModel,
		"areAuthFieldsInModel": func(cs []*Column, cl []*JWTField) bool {
			isInModel := false
			for _, c := range cl {
				if isAuthFieldInModel(cs, c) {
					isInModel = true
				}
			}
			return isInModel
		},
		"isUnique": func(c *Column) bool {
			for _, attr := range c.Attributes {
				if attr.Name == "unique" {
					return true
				}
			}

			return false
		},
		"isNotFirstUnique": func(cs []*Column, cl *Column) bool {
			count := 0
			for _, c := range cs {
				if c.SnakeName == cl.SnakeName {
					return count > 0
				}
				for _, attr := range c.Attributes {
					if attr.Name == "unique" {
						count++
					}
				}
			}
			return false
		},
		"isNullableForInput": func(c *Column) bool {
			if c.Type.IsNullable {
				return true
			}

			for _, attr := range c.Attributes {
				if attr.Name == "default" {
					return true
				}
			}

			return false
		},
		"isNullableDBType": func(c *Column) bool {
			if c.Type.IsNullable {
				return true
			}

			// for _, attr := range c.Attributes {
			// 	if attr.Name == "default" {
			// 		return false
			// 	}
			// }

			return false
		},
		"isInJwtField": isInJwtField,
		"isNotInJwtField": func(snakeName string, fs []*JWTField) bool {
			return !isInJwtField(snakeName, fs)
		},
		"isAuthTable": func(m *CreateAndUpdateDataModel) bool {
			return m.IsAuthUser || m.IsAuthOrganization || m.IsAuthOrganizationUser
		},
		"getAuthQueryInBoolAnd": func(f []*JWTField) string {
			var s []string
			for _, ff := range f {
				s = append(s, "!"+firstToLower(ff.TableCamelName)+"InBool")
			}
			return strings.Join(s, "&&")
		},
		"isRelationWithoutID": func(c *Column) bool {
			return c.IsRelation && !strings.HasSuffix(c.SnakeName, "_id")
		},
	}).Parse(c.Template)
	if err != nil {
		return "", fmt.Errorf("parse: %v", err)
	}

	var content bytes.Buffer
	err = tpl.Execute(&content, c.Data)
	if err != nil {
		return "", fmt.Errorf("execute: %v", err)
	}

	if shouldFormat {
		contentBytes := content.Bytes()
		formattedContent, err := format.Source(contentBytes)
		if err != nil {
			return string(contentBytes), fmt.Errorf("formatting: %v", err)
		}

		return string(formattedContent), nil
	}

	return content.String(), nil
}

func isInJwtField(snakeName string, fs []*JWTField) bool {
	for _, ff := range fs {
		if strings.EqualFold(ff.SnakeName, snakeName) {
			return true
		}
	}

	return false
}

func isAuthFieldInModel(cs []*Column, cl *JWTField) bool {
	for _, c := range cs {
		if c.IsRelation && strings.EqualFold(strings.TrimSuffix(c.SnakeName, "_id"), strings.TrimSuffix(cl.SnakeName, "_id")) {
			return true
		}
	}
	return false
}

func getRelations(f *ast.File, m *Model, pkgName string, models []*Model) []*ModelTemplateRelation {
	var relations []*ModelTemplateRelation
	lowerName := firstToLower(m.CamelName)

	ast.Inspect(f, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				switch typeSpec.Name.Name {
				// case m.CamelName:
				// 	// combineStructs = fmt.Sprintf("type %s struct {\n", typeSpec.Name.Name)

				// 	for _, field := range structType.Fields.List {
				// 		if field.Names[0].Name == "R" || field.Names[0].Name == "L" {
				// 			continue
				// 		}
				// 		fieldType, err := fieldTypeToString(field.Type)
				// 		if err != nil {
				// 			panic(err)
				// 		}

				// 		if strings.HasPrefix(fieldType, "time.") {
				// 			imports = addImport(imports, "time")
				// 		}

				// 		if strings.HasPrefix(fieldType, "null.") {
				// 			imports = addImport(imports, "github.com/volatiletech/null/v8")
				// 		}

				// 		// combineStructs += fmt.Sprintf("\t%s %s", field.Names[0].Name, fieldType)
				// 		// if field.Tag != nil {
				// 		// 	combineStructs += fmt.Sprintf(" %s\n", field.Tag.Value)
				// 		// } else {
				// 		// 	combineStructs += "\n"
				// 		// }
				// 	}

				case lowerName + "R":
					for _, field := range structType.Fields.List {
						fieldType, err := fieldTypeToString(pkgName+".", field.Type)
						if err != nil {
							panic(err)
						}

						singularName := field.Names[0].Name
						if isPlural(singularName) {
							singularName = singularize(field.Names[0].Name)
						}

						if strings.EqualFold(singularName, "oauth2") {
							// TODO: maybe hide @hide models here
							continue
						}

						var columns []*Column
						for _, m := range models {
							if m.CamelName == singularName {
								columns = m.Columns
							}
						}

						if strings.HasPrefix(fieldType, "*") {
							fieldType = strings.Replace(fieldType, "*"+pkgName+".", "*", 1)
						} else {
							fieldType = strings.Replace(fieldType, pkgName+".", "", 1)
						}

						r := ModelTemplateRelation{
							Name:         field.Names[0].Name,
							SingularName: singularName,
							Type:         fieldType,
							Columns:      columns,
							IsArray:      strings.Contains(fieldType, "Slice"),
						}

						if field.Tag != nil {
							r.Tag = field.Tag.Value
						}

						relations = append(relations, &r)
					}
				}
			}

		}

		return true
	})

	return relations
}

func everyRouteIsProtected(r []string) bool {
	allRoutes := []string{"LIST", "BYID", "CREATE", "UPDATE", "DELETE"}

	for _, v := range allRoutes {
		if !stringArrayContains(r, v) {
			return false
		}
	}

	return true

}

func stringArrayContains(arr []string, target string) bool {
	for _, s := range arr {
		if strings.EqualFold(s, target) {
			return true
		}
	}
	return false
}
