package functions

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ParseGoSQLFile(fileName string, sqlType string) (string, []*Model) {
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
	sqlTypeMatch := sqlTypeCompiled.FindString(strings.Split(str, "\n")[0])
	if len(sqlTypeMatch) <= 1 && sqlType == "" {
		panic("Invalid sql type")
	}
	if sqlType == "" {
		sqlType = sqlTypeMatch[1:]
	}

	re := regexp.MustCompile(modelRegexp3)
	match := re.FindAllStringSubmatch(str, -1)

	var hasAuthUser bool
	var hasAuthOrganization bool
	var hasAuthOrganizationUser bool

	for _, m := range match {
		extraAttributeLine := m[3]

		if len(m) <= 1 {
			panic("Invalid model definition")
		}

		name := camelToSnake(m[1])
		body := strings.TrimSpace(m[2])

		attrRegex := regexp.MustCompile(`@?(\w+)\(([^)]*)\)|@(\w+)`)
		attrMatch := attrRegex.FindAllStringSubmatch(extraAttributeLine, -1)

		protectedRoutes := []string{}
		isAuthUser := false
		isAuthOrganization := false
		isAuthOrganizationUser := false

		for _, match := range attrMatch {
			if match[1] != "" {
				// fmt.Printf("%s: %v\n", match[1], match[2])

				if match[1] == "protected" {
					x := strings.Split(match[2], ",")
					for _, y := range x {
						protectedRoutes = append(protectedRoutes, strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(y, `"`, ""), `'`, "")))
					}
				}

			} else {
				// fmt.Printf("%s\n", match[3])
				if match[3] == "protected" {
					protectedRoutes = []string{"LIST", "BYID", "CREATE", "UPDATE", "DELETE"}
				}

				if strings.EqualFold(match[3], "authUser") {
					isAuthUser = true
					hasAuthUser = true
				}

				if strings.EqualFold(match[3], "authOrganization") {
					isAuthOrganization = true
					hasAuthOrganization = true
				}

				if strings.EqualFold(match[3], "authOrganizationUser") {
					isAuthOrganizationUser = true
					hasAuthOrganizationUser = true
				}
			}
		}

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
			t.TypescriptName = getTypescriptType(t.Name)

			columns = append(columns, &Column{
				SnakeName:  camelToSnake(splittedLine[0]),
				CamelName:  snakeToCamel(splittedLine[0]),
				Type:       t,
				Attributes: attributes,
				IsRelation: isRelation(t.Name),
				Expose:     shouldExpose(attributes),
			})
		}

		if isAuthUser {
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
				panic("AuthUser model must have a unique email column")
			}

			if !isValidPasswordColumn {
				panic("AuthUser model must have a password column")
			}
		}

		m := &Model{
			SnakeName:              name,
			CamelName:              snakeToCamel(name),
			Columns:                columns,
			IsAuthRequired:         false, // @deprecated
			IsAuthUser:             isAuthUser,
			ProtectedRoutes:        protectedRoutes,
			IsAuthOrganization:     isAuthOrganization,
			IsAuthOrganizationUser: isAuthOrganizationUser,
		}

		models = append(models, m)

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
	}

	if hasAuthOrganization && (!hasAuthUser || !hasAuthOrganizationUser) {
		panic("If you have an organization model you must have an authUser and authOrganizationUser table")
	}

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

func getTypescriptType(t string) string {
	switch t {
	case "string":
		return "string"
	case "text":
		return "string"
	case "uuid":
		return "string"
	case "int":
		return "number"
	case "bool":
		return "boolean"
	case "dateTime":
		return "Date"
	case "uint":
		return "number"
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

// 2056 - 1303 = 753
