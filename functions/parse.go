package functions

import (
	"os"
	"regexp"
	"strconv"
	"strings"
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

	dat := make([]byte, 1024)
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

		name := match[1]
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
					re := regexp.MustCompile(attributeWithValueRegexp)
					match := re.FindStringSubmatch(extraAttribute)
					name := match[1]
					value := match[2]
					attributes = append(attributes, &Attribute{
						Name:     name,
						Value:    value,
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
				Name:       splittedLine[0],
				Type:       t,
				Attributes: attributes,
			})
		}

		models = append(models, &Model{
			Name:    name,
			Columns: columns,
		})

	}

	return sqlType, models

}
