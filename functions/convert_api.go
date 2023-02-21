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

	for _, m := range models {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fmt.Sprintf("%s/%s.go", config.Output, m.SnakeName), nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("make sure you run sqlboiler first")
		}

		var imports []string
		imports = addImport(imports, moduleName+"/"+config.Output) // sql boiler models
		relations := getRelations(f, m, pkgName, models)

		if err := populateTemplate("templates/model.gotpl", outputDir+"/"+m.SnakeName+".go", ModelTemplateData{
			PackageName: "am",
			CamelName:   m.CamelName,
			Imports:     imports,
			Relations:   relations,
		}); err != nil {
			return err
		}
	}

	return nil

}

func (c *GoSQLConfig) ConvertApiControllers(models []*Model) error {
	outputDir := c.ControllerOutputDir
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

	if err := populateTemplate("./templates/helpers.gotpl", outputDir+"/generated_helpers.go", TemplateData{PackageName: "controllers"}); err != nil {
		return err
	}
	if err := populateTemplate("templates/queries.gotpl", outputDir+"/generated_queries.go", GeneralTemplateData{PackageName: "controllers", Controllers: models}); err != nil {
		return err
	}
	if err := populateTemplate("templates/filters.gotpl", outputDir+"/generated_filters.go", GeneralTemplateData{PackageName: "controllers", Controllers: models}); err != nil {
		return err
	}
	if err := populateTemplate("templates/client.gotpl", outputDir+"/generated_client.go", GeneralTemplateData{PackageName: "controllers", Controllers: models}); err != nil {
		return err
	}
	if err := populateTemplate("templates/columns.gotpl", outputDir+"/generated_columns.go", GeneralTemplateData{PackageName: "controllers", Controllers: models}); err != nil {
		return err
	}

	if err := populateTemplate("templates/routes.gotpl", outputDir+"/generated_routes.go", GeneralTemplateData{PackageName: "controllers", Controllers: models}); err != nil {
		return err
	}

	if err := populateTemplate("templates/orders.gotpl", outputDir+"/generated_orders.go", GeneralTemplateData{PackageName: "controllers", Controllers: models}); err != nil {
		return err
	}

	// Build the controllers
	var modelWithRelations []*ModelWithRelations
	var createAndUpdateData []*CreateAndUpdateDataModel

	var modelImports []string
	modelImports = addImport(modelImports, moduleName+"/"+config.Output)

	for _, m := range models {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fmt.Sprintf("%s/%s.go", config.Output, m.SnakeName), nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("make sure you run sqlboiler first")
		}
		relations := getRelations(f, m, pkgName, models)
		modelWithRelations = append(modelWithRelations, &ModelWithRelations{Model: m, Relations: relations})

		createColumns, updateColumns, mImports := getCreateAndUpdateColumns(m)

		if err := populateTemplate("templates/controller.gotpl", outputDir+"/generated_"+m.SnakeName+"_controller.go", ControllerTemplateData{PackageName: "controllers", CamelName: m.CamelName, Imports: imports, CreateColumns: createColumns, UpdateColumns: updateColumns}); err != nil {
			return err
		}

		modelImports = append(modelImports, mImports...)
		createAndUpdateData = append(createAndUpdateData, &CreateAndUpdateDataModel{
			SnakeName:     m.SnakeName,
			CamelName:     m.CamelName,
			CreateColumns: createColumns,
			UpdateColumns: updateColumns,
		})
	}

	if err := populateTemplate("templates/relations.gotpl", outputDir+"/generated_relations.go", SelectTemplateData{PackageName: "controllers", Controllers: modelWithRelations}); err != nil {
		return err
	}

	if err := populateTemplate("templates/bodies.gotpl", outputDir+"/generated_bodies.go", CreateAndUpdateData{PackageName: "controllers", Controllers: createAndUpdateData, Imports: modelImports}); err != nil {
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

	template, err := parseTemplate(&TemplateConfig{Template: string(content), Data: data})
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

type TemplateConfig struct {
	Template string
	Data     interface{}
}

func parseTemplate(c *TemplateConfig) (string, error) {
	tpl, err := template.New("").Funcs(template.FuncMap{
		"toSnake":      camelToSnake,
		"pluralize":    pluralize,
		"firstToLower": firstToLower,
		"isFalse": func(a bool) bool {
			return !a
		},
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
		"isNullable": func(c *Column) bool {
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
	}).Parse(c.Template)
	if err != nil {
		return "", fmt.Errorf("parse: %v", err)
	}

	var content bytes.Buffer
	err = tpl.Execute(&content, c.Data)
	if err != nil {
		return "", fmt.Errorf("execute: %v", err)
	}

	contentBytes := content.Bytes()
	formattedContent, err := format.Source(contentBytes)
	if err != nil {
		return string(contentBytes), fmt.Errorf("formatting: %v", err)
	}

	return string(formattedContent), nil
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

						var columns []*Column
						for _, m := range models {
							if m.CamelName == singularName {
								columns = m.Columns
							}
						}

						r := ModelTemplateRelation{
							Name:         field.Names[0].Name,
							SingularName: singularName,
							Type:         fieldType,
							Columns:      columns,
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
