package functions

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

func (c *GoSQLConfig) ConvertTypes(models []*Model) error {
	outputDir := "types"

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

	pkgName := config.PkgName

	var modelWithRelationsAndIds []*Model
	for _, model := range models {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fmt.Sprintf("%s/%s.go", config.Output, model.SnakeName), nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("make sure you run sqlboiler first")
		}

		relations := getRelations(f, model, pkgName, models)

		var customColumns []*Column
		for _, column := range model.Columns {

			if column.IsRelation {
				customColumns = append(customColumns, &Column{
					SnakeName:  column.SnakeName + "_id",
					CamelName:  column.CamelName + "Id",
					Type:       &Type{Name: "int", GoTypeName: "int", TypescriptName: "number"},
					Attributes: column.Attributes,
					IsRelation: column.IsRelation,
					Expose:     column.Expose,
				})
			} else {
				customColumns = append(customColumns, column)
			}
		}
		for _, relation := range relations {
			customColumns = append(customColumns, &Column{
				SnakeName:  camelToSnake(relation.Name),
				CamelName:  relation.Name,
				Type:       getRelationType(relation),
				Attributes: []*Attribute{},
				IsRelation: true,
				Expose:     true,
			})
		}

		modelWithRelationsAndIds = append(modelWithRelationsAndIds, &Model{
			SnakeName:              model.SnakeName,
			CamelName:              model.CamelName,
			Columns:                customColumns,
			IsAuthRequired:         model.IsAuthRequired,
			IsAuthUser:             model.IsAuthUser,
			ProtectedRoutes:        model.ProtectedRoutes,
			IsAuthOrganization:     model.IsAuthOrganization,
			IsAuthOrganizationUser: model.IsAuthOrganizationUser,
		})
	}

	if err := populateTemplate("templates/typescript/types.gotpl", outputDir+"/types.ts", TypescriptTypesTemplateData{
		Controllers: modelWithRelationsAndIds,
	}); err != nil {
		return err
	}

	if err := populateTemplate("templates/typescript/requests.gotpl", outputDir+"/requests.ts", TypescriptTypesTemplateData{
		Controllers: modelWithRelationsAndIds,
	}); err != nil {
		return err
	}

	return nil
}

func getRelationType(relation *ModelTemplateRelation) *Type {
	if relation.IsArray {
		return &Type{Name: "[]" + relation.SingularName, GoTypeName: "[]" + relation.SingularName, TypescriptName: relation.SingularName + "[]"}
	}

	return &Type{Name: relation.SingularName, GoTypeName: relation.SingularName, TypescriptName: relation.SingularName}
}
