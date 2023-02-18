package functions

const (
	// body = `(?m)^\s{4}.*$`
	atSignRegexp              = `@([a-zA-Z]+)`
	modelRegexp               = `(?:^|\n)\s*([A-Z][a-zA-Z]*)\s*{((?:.|\n)*?)}`
	modelSplitRegexp          = `^\s*([A-Z][a-zA-Z]*)\s*{((?:.|\n)*?)}`
	attributeWithValueRegexp  = `^@([a-zA-Z]+)\(([^)]+)\)`
	attributeWithValueRegexp2 = `^@([a-zA-Z]+)\(([\s\S]+)\)`
)

type Model struct {
	SnakeName string
	CamelName string
	Columns   []*Column
}

type ModelWithRelations struct {
	*Model
	Relations []*ModelTemplateRelation
}

type Column struct {
	SnakeName  string
	CamelName  string
	Type       *Type
	Attributes []*Attribute
}

type Type struct {
	Name                   string
	IsNullable             bool
	HasDifferentCharLength bool
	CharLength             int
}

type Attribute struct {
	Name     string
	Value    string
	HasValue bool
}

type TemplateData struct {
	PackageName string
}

type GeneralTemplateData struct {
	PackageName string
	Controllers []*Model
}

type SelectTemplateData struct {
	PackageName string
	Controllers []*ModelWithRelations
}

type ControllerTemplateData struct {
	PackageName string
	CamelName   string
	Imports     []string
}

type ModelTemplateRelation struct {
	Name string
	Type string
	Tag  string
}

type ModelTemplateData struct {
	PackageName string
	Imports     []string
	CamelName   string
	Relations   []*ModelTemplateRelation
}
