package functions

import (
	"fmt"
	"os"
	"strings"
)

func (c *GoSQLConfig) ConvertToSql(fileName, t string, models []*Model, existingModels []*Model) error {
	if t != "postgresql" {
		return fmt.Errorf("sql type %s not supported", t)
	}
	fileName = strings.TrimSuffix(fileName, ".gosql")

	file, err := os.Create(c.MigrationDir + "/" + fileName + ".up.sql")
	if err != nil {
		return err
	}
	defer file.Close()

	for _, m := range models {
		var constraints []string
		// var relations []string
		var indexes []string
		var ak = 1
		var fk = 1
		var idx = 1
		var requiresUuidExtension bool

		tableName := strings.ToLower(m.SnakeName)
		file.WriteString("CREATE TABLE " + tableName + " (\n")
		for _, c := range m.Columns {
			var line string

			t, isRelation, rUe := convertType(tableName, c, func(constraintType string) {
				var constraint string
				if constraintType == "ID" {
					constraint = fmt.Sprintf("CONSTRAINT %v_pk PRIMARY KEY (%v)", tableName, c.SnakeName)
				}
				if constraintType == "UNIQUE" {
					constraint = fmt.Sprintf("CONSTRAINT %v_ak_%v UNIQUE (%v) NOT DEFERRABLE INITIALLY IMMEDIATE", tableName, ak, c.SnakeName)
					ak++
				}
				constraints = append(constraints, constraint)
			})

			if rUe {
				requiresUuidExtension = true
			}

			if isRelation {

				// Check if exists in models
				var exists bool
				var relationType string
				for _, m := range existingModels {
					if strings.EqualFold(m.SnakeName, t) {
						exists = true
						for _, co := range m.Columns {
							if strings.EqualFold(co.SnakeName, "id") {
								relationType, _, _ = convertType(tableName, co, func(constraintType string) {})
								if relationType == "SERIAL" {
									relationType = "INTEGER"
								}
							}
						}
					}
				}

				if !exists || relationType == "" {
					return fmt.Errorf("Relation " + t + " does not exist")
				}

				snakeColumnName := c.SnakeName + "_id"
				line = fmt.Sprintf("\t\t%v %v", snakeColumnName, relationType)

				if c.Type.IsNullable {
					line += " NULL"
				} else {
					line += " NOT NULL"
				}

				constraints = append(constraints, fmt.Sprintf("CONSTRAINT %v_%v_fk_%v FOREIGN KEY (%v) REFERENCES %v (id) NOT DEFERRABLE INITIALLY IMMEDIATE", tableName, snakeColumnName, fk, snakeColumnName, t))
				fk++

				for _, a := range c.Attributes {
					if a.Name == "index" {
						indexes = append(indexes, fmt.Sprintf("CREATE INDEX %v_idx_%v ON %v (%v)", tableName, idx, tableName, snakeColumnName))
						idx++
					}
				}

			} else {
				line = fmt.Sprintf("\t\t%v %v", c.SnakeName, t)
				if c.Type.IsNullable {
					line += " NULL"
				} else {
					line += " NOT NULL"
				}

				for _, a := range c.Attributes {
					if a.Name == "default" && a.HasValue {
						if a.Value == "autoincrement" {
							continue
						}
						if a.Value == "now" {
							line += " DEFAULT NOW()"
							continue
						}

						line += " DEFAULT " + a.Value
					}

					if a.Name == "index" {
						indexes = append(indexes, fmt.Sprintf("CREATE INDEX %v_idx_%v ON %v (%v)", tableName, idx, tableName, c.SnakeName))
						idx++
					}
				}
			}

			file.WriteString(line)
			file.WriteString(",\n")
		}

		for i, constraint := range constraints {
			line := fmt.Sprintf("\t\t%v", constraint)
			if i != len(constraints)-1 {
				line += ","
			}
			file.WriteString(line + "\n")
		}

		file.WriteString(");\n")

		if requiresUuidExtension {
			file.WriteString("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";\n")
		}

		for _, index := range indexes {
			line := fmt.Sprintf("%v;", index)
			file.WriteString(line + "\n")
		}

		file.WriteString("\n")
	}

	return nil

}

func convertType(tableName string, c *Column, cb func(constraint string)) (typeName string, isRelation bool, requiresUuidExtension bool) {
	for _, a := range c.Attributes {
		if a.Name == "default" && a.HasValue && a.Value == "autoincrement" && c.Type.Name == "int" {
			cb("ID")
			return "SERIAL", false, false
		}

		if a.Name == "default" && a.HasValue && a.Value == "uuid" && c.Type.Name == "uuid" {
			cb("ID")
			return "UUID", false, true
		}

		if a.Name == "unique" {
			cb("UNIQUE")
		}
	}

	if c.IsRelation {
		return strings.ToLower(camelToSnake(c.Type.Name)), true, false
	}

	switch c.Type.Name {
	case "string":
		return fmt.Sprintf("VARCHAR(%v)", c.Type.CharLength), false, false
	case "text":
		return "TEXT", false, false
	case "int":
		return "INTEGER", false, false
	case "bool":
		return "BOOLEAN", false, false
	case "dateTime":
		return "TIMESTAMP", false, false
	case "uuid":
		return "UUID", false, false
	case "uint":
		return "UINT", false, false
	default:
		return strings.ToUpper(c.Type.Name), false, false
	}
}
