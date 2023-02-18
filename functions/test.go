package functions

// func parse{{ .CamelName }}FilterToMods(query *{{.CamelName}}ListQuery) ([]qm.QueryMod, error) {
// 	var queryMods []qm.QueryMod

// 	queryMods = append(queryMods, qm.Offset(query.Page))
// 	queryMods = append(queryMods, qm.Limit(query.Limit))

// 	if query.Filter != nil {
// 		{{ range $column := .Columns }}
// 			if query.Filter.{{ .CamelName }} != nil {
// 				if query.Filter.{{ .CamelName }}.Equals != nil {
// 					queryMods = append(queryMods, qm.Where("{{ .SnakeName }} = ?", *query.Filter.{{ .CamelName }}.Equals))
// 				}

// 				if query.Filter.{{ .CamelName }}.NotEquals != nil {
// 					queryMods = append(queryMods, qm.Where("{{ .SnakeName }} != ?", *query.Filter.{{ .CamelName }}.NotEquals))
// 				}

// 				if query.Filter.{{ .CamelName }}.IsIn != nil {
// 					queryMods = append(queryMods, qm.WhereIn("{{ .SnakeName }} IN ?", *query.Filter.{{ .CamelName }}.IsIn))
// 				}

// 				if query.Filter.{{ .CamelName }}.IsNotIn != nil {
// 					queryMods = append(queryMods, qm.WhereIn("{{ .SnakeName }} NOT IN ?", *query.Filter.{{ .CamelName }}.IsNotIn))
// 				}

// 				if query.Filter.{{ .CamelName }}.IsLessThan != nil {
// 					queryMods = append(queryMods, qm.Where("{{ .SnakeName }} < ?", *query.Filter.{{ .CamelName }}.IsLessThan))
// 				}

// 				if query.Filter.{{ .CamelName }}.IsLessThanOrEqual != nil {
// 					queryMods = append(queryMods, qm.Where("{{ .SnakeName }} <= ?", *query.Filter.{{ .CamelName }}.IsLessThanOrEqual))
// 				}

// 				if query.Filter.{{ .CamelName }}.IsGreaterThan != nil {
// 					queryMods = append(queryMods, qm.Where("{{ .SnakeName }} > ?", *query.Filter.{{ .CamelName }}.IsGreaterThan))
// 				}

// 				if query.Filter.{{ .CamelName }}.IsGreaterThanOrEqual != nil {
// 					queryMods = append(queryMods, qm.Where("{{ .SnakeName }} >= ?", *query.Filter.{{ .CamelName }}.IsGreaterThanOrEqual))
// 				}
// 			}
// 		{{- end }}
// 	}

// 	return queryMods, nil
// }

// func {{ .CamelName }}FilterToMods(context *gin.Context) ([]qm.QueryMod, error) {
// 	query, err := parse{{ .CamelName }}ListQueryParameters(context)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return parse{{ .CamelName }}FilterToMods(query)
// }
