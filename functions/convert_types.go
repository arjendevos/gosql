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

	var authQueryFields []*JWTField
	var modelWithRelationIds []*Model
	var modelWithRelationsAndRelationIds []*ModelWithRelations

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
					SnakeName:    column.SnakeName + "_id",
					CamelName:    column.CamelName + "Id",
					Type:         &Type{Name: "int", GoTypeName: "int", TypescriptName: "number"},
					Attributes:   column.Attributes,
					IsRelation:   column.IsRelation,
					Expose:       column.Expose,
					DatabaseName: column.DatabaseName,
				})
			} else {
				customColumns = append(customColumns, column)
			}
		}

		for _, relation := range relations {
			customColumns = append(customColumns, &Column{
				SnakeName:    camelToSnake(relation.Name),
				CamelName:    relation.Name,
				Type:         getRelationType(relation),
				Attributes:   []*Attribute{},
				IsRelation:   true,
				Expose:       true,
				DatabaseName: relation.DatabaseName,
			})
		}

		theModelWithRelationIds := Model{
			SnakeName:              model.SnakeName,
			CamelName:              model.CamelName,
			Columns:                customColumns,
			IsAuthRequired:         model.IsAuthRequired,
			IsAuthUser:             model.IsAuthUser,
			ProtectedRoutes:        model.ProtectedRoutes,
			IsAuthOrganization:     model.IsAuthOrganization,
			IsAuthOrganizationUser: model.IsAuthOrganizationUser,
			HideRoutes:             model.HideRoutes,
			Hide:                   model.Hide,
			Oauth2:                 model.Oauth2,
		}

		modelWithRelationIds = append(modelWithRelationIds, &theModelWithRelationIds)
		modelWithRelationsAndRelationIds = append(modelWithRelationsAndRelationIds, &ModelWithRelations{Model: &theModelWithRelationIds, Relations: relations})

		if model.IsAuthUser {
			for _, c := range model.Columns {
				if c.SnakeName == "id" {
					authQueryFields = append(authQueryFields, &JWTField{
						CamelName:                   model.CamelName + c.CamelName,
						SnakeName:                   model.SnakeName + "_" + c.SnakeName,
						GoType:                      c.Type.GoTypeName,
						NormalName:                  c.CamelName,
						TableCamelName:              model.CamelName,
						TableSnakeName:              model.SnakeName,
						IsFromUserTable:             true,
						IsFromOrganizationTable:     false,
						IsFromOrganizationUserTable: false,
					})
				}
			}
		}

		if model.IsAuthOrganization {
			for _, c := range model.Columns {
				if c.SnakeName == "id" {
					authQueryFields = append(authQueryFields, &JWTField{
						CamelName:                   model.CamelName + c.CamelName,
						SnakeName:                   model.SnakeName + "_" + c.SnakeName,
						GoType:                      c.Type.GoTypeName,
						NormalName:                  c.CamelName,
						TableCamelName:              model.CamelName,
						TableSnakeName:              model.SnakeName,
						IsFromUserTable:             false,
						IsFromOrganizationTable:     true,
						IsFromOrganizationUserTable: false,
					})
				}
			}
		}
	}

	if err := populateTemplate("templates/typescript/types.gotpl", outputDir+"/types.ts", TypescriptTypesTemplateData{
		Controllers: modelWithRelationIds,
	}); err != nil {
		return err
	}

	if err := populateTemplate("templates/typescript/query.gotpl", outputDir+"/query.ts", TypescriptTypesRelationsTemplateData{
		Controllers:           modelWithRelationsAndRelationIds,
		AuthFields:            authQueryFields,
		HasMultipleAuthFields: len(authQueryFields) > 0,
	}); err != nil {
		return err
	}

	if err := populateTemplate("templates/typescript/requests.gotpl", outputDir+"/requests.ts", TypescriptTypesTemplateData{
		Controllers: modelWithRelationIds,
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
