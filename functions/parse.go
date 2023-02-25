package functions

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func parseGoSQLFile(fileName string) (string, []*Model) {
	var sqlType string
	var models []*Model

	fileName = strings.TrimSuffix(fileName, ".gosql")
	file, err := os.Open(fileName + ".gosql")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	size, err := file.Stat()
	if err != nil {
		panic(err)
	}

	dat := make([]byte, size.Size())
	_, err = file.Read(dat)
	if err != nil {
		panic(err)
	}
	str := string(dat)

	sqlTypeCompiled := regexp.MustCompile(atSignRegexp)
	sqlTypeMatch := sqlTypeCompiled.FindString(str)
	sqlType = sqlTypeMatch[1:]

	re := regexp.MustCompile(modelRegexp2)
	match := re.FindAllStringSubmatch(str, -1)

	var hasAuth bool
	var authModel *Model

	for _, m := range match {
		atStrings := m[3:]
		// re := regexp.MustCompile(modelSplitRegexp)
		// match := re.FindStringSubmatch(model)

		if len(match) <= 1 {
			panic("Invalid model definition")
		}

		name := camelToSnake(m[1])
		body := strings.TrimSpace(m[2])

		var columns []*Column

		bodyArray := strings.Split(body, "\n")

		for _, line := range bodyArray {
			splittedLine := strings.Split(strings.TrimSpace(line), " ")
			if len(splittedLine) <= 1 {
				panic("Invalid column definition")
			}

			var attributes []*Attribute
			extraAttributesArray := splittedLine[2:]
			for _, extraAttribute := range extraAttributesArray {
				if !strings.Contains(extraAttribute, "@") {
					panic("Invalid attribute definition")
				}

				if strings.Contains(extraAttribute, "(") && strings.Contains(extraAttribute, ")") {
					re := regexp.MustCompile(attributeWithValueRegexp2)
					match := re.FindStringSubmatch(extraAttribute)

					name := match[1]
					value := match[2]
					attributes = append(attributes, &Attribute{
						Name:     name,
						Value:    strings.ReplaceAll(value, `"`, "'"),
						HasValue: true,
					})

				} else {
					name := strings.ReplaceAll(extraAttribute, "@", "")
					attributes = append(attributes, &Attribute{
						Name:     name,
						HasValue: false,
					})
				}

			}

			t := &Type{
				Name:                   splittedLine[1],
				IsNullable:             false,
				HasDifferentCharLength: false,
				CharLength:             255,
			}

			if strings.HasSuffix(t.Name, "?") {
				t.Name = strings.TrimSuffix(t.Name, "?")
				t.IsNullable = true
			}

			if strings.Contains(t.Name, "(") && strings.Contains(t.Name, ")") {
				t.HasDifferentCharLength = true
				re := regexp.MustCompile(`\((.*?)\)`)
				match := re.FindStringSubmatch(t.Name)
				t.CharLength, err = strconv.Atoi(match[1])
				if err != nil {
					panic(err)
				}
				t.Name = strings.ReplaceAll(t.Name, match[0], "")
			}

			t.GoTypeName = getGotype(t.Name)

			columns = append(columns, &Column{
				SnakeName:  camelToSnake(splittedLine[0]),
				CamelName:  snakeToCamel(splittedLine[0]),
				Type:       t,
				Attributes: attributes,
				IsRelation: isRelation(t.Name),
				Expose:     shouldExpose(attributes),
			})
		}

		isAuthRequired := true
		var isAuthUser, isAuthOrg, isAuthOrgLink bool

		for _, a := range atStrings {
			if strings.Contains(a, "authUser") {
				isAuthUser = true

				var isValidEmailColumn bool
				var isValidPasswordColumn bool
				for _, c := range columns {
					if c.SnakeName == "email" && c.Type.Name == "string" {
						for _, a := range c.Attributes {
							if a.Name == "unique" {
								isValidEmailColumn = true
							}
						}
					}

					if c.SnakeName == "password" && c.Type.Name == "string" {
						isValidPasswordColumn = true
					}
				}

				if !isValidEmailColumn {
					panic("Auth user model must have a unique email column")
				}

				if !isValidPasswordColumn {
					panic("Auth user model must have a password column")
				}
			}

			if strings.Contains(a, "authOrg") {
				isAuthOrg = true
			}

			if strings.Contains(a, "authOrgLink") {
				isAuthOrgLink = true
			}

			if strings.Contains(a, "noAuth") {
				isAuthRequired = false
			}
		}

		m := &Model{
			SnakeName:      name,
			CamelName:      snakeToCamel(name),
			Columns:        columns,
			IsAuthRequired: isAuthRequired,
			IsAuthUser:     isAuthUser,
			IsAuthOrg:      isAuthOrg,
			IsAuthOrgLink:  isAuthOrgLink,
		}

		models = append(models, m)

		if m.IsAuthUser {
			hasAuth = true
			authModel = m
		}
	}

	if hasAuth {
		for _, m := range models {
			if m.IsAuthUser {
				continue
			}

			// var idType *Type
			// for _, c := range m.Columns {
			// 	if c.SnakeName == "id" {
			// 		idType = c.Type
			// 	}
			// }

			var hasAuthUser bool
			for _, c := range m.Columns {
				if c.IsRelation && strings.EqualFold(c.Type.Name, authModel.CamelName) {
					hasAuthUser = true
				}
			}

			if hasAuthUser {
				continue
			}

			// m.Columns = append(m.Columns, &Column{
			// 	SnakeName: camelToSnake(authModel.CamelName),
			// 	CamelName: snakeToCamel(authModel.SnakeName),
			// 	Type: &Type{
			// 		Name:                   firstToLower(authModel.CamelName),
			// 		GoTypeName:             authModel.CamelName,
			// 		IsNullable:             false,
			// 		HasDifferentCharLength: false,
			// 		CharLength:             255,
			// 	},
			// 	Attributes: []*Attribute{
			// 		{
			// 			Name:     "index",
			// 			Value:    "",
			// 			HasValue: false,
			// 		},
			// 	},
			// 	IsRelation: true,
			// })
		}
	}

	// for _, m := range models {
	// 	fmt.Println(m.Name)
	// 	for _, c := range m.Columns {
	// 		fmt.Println(c.Name)
	// 		fmt.Println("  ", c.Type.Name)
	// 		fmt.Println("  ", c.Type.IsNullable)
	// 		fmt.Println("  ", c.Type.HasDifferentCharLength)
	// 		fmt.Println("  ", c.Type.CharLength)

	// 		for _, a := range c.Attributes {
	// 			fmt.Println("    ", a.Name, a.Value, a.HasValue)
	// 		}
	// 	}
	// }

	return sqlType, models

}

func snakeToCamel(s string) string {
	s = camelToSnake(s)

	if strings.Contains(strings.ToLower(s), "id") {
		i := strings.Index(strings.ToLower(s), "id")
		if len(s) <= i+2 {
			s = strings.ReplaceAll(strings.ToLower(s), "id", "ID")
		}
	}

	if strings.Contains(strings.ToLower(s), "url") {
		s = strings.ReplaceAll(strings.ToLower(s), "url", "URL")
	}

	if strings.Contains(strings.ToLower(s), "csfr") {
		s = strings.ReplaceAll(strings.ToLower(s), "csfr", "CSFR")
	}

	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Title(s)
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func camelToSnake(input string) string {
	var output bytes.Buffer
	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				output.WriteRune('_')
			}
			output.WriteRune(unicode.ToLower(r))
		} else {
			output.WriteRune(r)
		}
	}
	return output.String()
}

func getGotype(t string) string {
	switch t {
	case "string":
		return "string"
	case "text":
		return "string"
	case "uuid":
		return "string"
	case "int":
		return "int"
	case "bool":
		return "bool"
	case "dateTime":
		return "time.Time"
	case "uint":
		return "uint"
	default:
		return t
	}
}

func isRelation(s string) bool {
	return unicode.IsUpper(rune(s[0]))
}

func shouldExpose(at []*Attribute) bool {
	for _, a := range at {
		if a.Name == "hide" {
			return false
		}
	}

	return true
}
