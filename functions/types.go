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
	Name    string
	Columns []*Column
}

type Column struct {
	Name       string
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
