package functions

import "os"

func (c *GoSQLConfig) ConvertTypes(models []*Model) error {
	outputDir := "types"

	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	var modelWithRelationsAndIds []*Model
	for _, model := range models {
		var customColumns []*Column
		for _, column := range model.Columns {
			customColumns = append(customColumns, column)

			if column.IsRelation {
				customColumns = append(customColumns, &Column{
					SnakeName:  column.SnakeName + "_id",
					CamelName:  column.CamelName + "Id",
					Type:       &Type{Name: "int", GoTypeName: "int", TypescriptName: "number"},
					Attributes: column.Attributes,
					IsRelation: column.IsRelation,
					Expose:     column.Expose,
				})
			}
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

	return nil
}
