package functions

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func ConvertPostgreSQL(fileName string, models []*Model) {
	fileName = strings.TrimSuffix(fileName, ".gosql")

	file, err := os.Create(fileName + ".sql")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, m := range models {
		var constraints []string
		// var relations []string
		var indexes []string
		var ak = 1
		var fk = 1
		var idx = 1

		tableName := strings.ToLower(m.Name)
		file.WriteString("CREATE TABLE " + tableName + " (\n")
		for _, c := range m.Columns {
			var line string

			t, isRelation := convertType(tableName, c, func(constraintType string) {
				var constraint string
				if constraintType == "ID" {
					constraint = fmt.Sprintf("CONSTRAINT %v_pk PRIMARY KEY (%v)", tableName, c.Name)
				}
				if constraintType == "UNQIQUE" {
					constraint = fmt.Sprintf("CONSTRAINT %v_ak_%v UNIQUE (%v) NOT DEFERRABLE INITIALLY IMMEDIATE", tableName, ak, c.Name)
					ak++
				}
				constraints = append(constraints, constraint)
			})

			if isRelation {
				columnName := c.Name + "_id"
				line = fmt.Sprintf("\t\t%v INT", columnName)

				if c.Type.IsNullable {
					line += " NULL"
				} else {
					line += " NOT NULL"
				}

				constraints = append(constraints, fmt.Sprintf("CONSTRAINT %v_%v_fk_%v FOREIGN KEY (%v) REFERENCES %v (id) NOT DEFERRABLE INITIALLY IMMEDIATE", tableName, columnName, fk, c.Name, t))
				fk++

				for _, a := range c.Attributes {
					if a.Name == "index" {
						indexes = append(indexes, fmt.Sprintf("CREATE INDEX %v_idx_%v ON %v (%v)", tableName, idx, tableName, columnName))
					}
				}

			} else {
				line = fmt.Sprintf("\t\t%v %v", c.Name, t)
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
						indexes = append(indexes, fmt.Sprintf("CREATE INDEX %v_idx_%v ON %v (%v)", tableName, idx, tableName, c.Name))
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

		for _, index := range indexes {
			line := fmt.Sprintf("%v;", index)
			file.WriteString(line + "\n")
		}

		file.WriteString("\n")
	}

}

func convertType(tableName string, c *Column, cb func(constraint string)) (typeName string, isRelation bool) {
	for _, a := range c.Attributes {
		if a.Name == "default" && a.HasValue && a.Value == "autoincrement" && c.Type.Name == "int" {
			cb("ID")
			return "SERIAL", false
		}

		if a.Name == "unique" {
			cb("UNQIQUE")
		}
	}

	if unicode.IsUpper(rune(c.Type.Name[0])) {
		return strings.ToLower(c.Type.Name), true
	}

	switch c.Type.Name {
	case "string":
		return fmt.Sprintf("VARCHAR(%v)", c.Type.CharLength), false
	case "text":
		return "TEXT", false
	case "int":
		return "INTEGER", false
	case "bool":
		return "BOOLEAN", false
	case "dateTime":
		return "TIMESTAMP", false
	default:
		return "UNKNOWN", false
	}
}
