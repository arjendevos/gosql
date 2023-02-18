package functions

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ParseGoSQLFile(fileName string) (string, []*Model) {
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

	re := regexp.MustCompile(modelRegexp)
	match := re.FindAllString(str, -1)

	for _, m := range match {
		re := regexp.MustCompile(modelSplitRegexp)
		match := re.FindStringSubmatch(m)

		if len(match) <= 1 {
			panic("Invalid model definition")
		}

		name := camelToSnake(match[1])
		body := strings.TrimSpace(match[2])

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

			columns = append(columns, &Column{
				SnakeName:  camelToSnake(splittedLine[0]),
				CamelName:  snakeToCamel(splittedLine[0]),
				Type:       t,
				Attributes: attributes,
			})
		}

		models = append(models, &Model{
			SnakeName: name,
			CamelName: snakeToCamel(name),
			Columns:   columns,
		})

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
