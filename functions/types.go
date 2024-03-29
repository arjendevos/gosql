package functions

const (
	// body = `(?m)^\s{4}.*$`
	atSignRegexp              = `@([a-zA-Z]+)`
	modelRegexp               = `(?:^|\n)\s*([A-Z][a-zA-Z]*)\s*{((?:.|\n)*?)}`
	modelRegexp2              = `(?:^|\n)\s*([A-Z][a-zA-Z]*)\s*{((?:.|\n)*?)}(?:\s*@(\w+))*`
	modelRegexp3              = `(?:^|\n)\s*([A-Z][a-zA-Z]*)\s*{((?:.|\n)*?)}(?:\s*(.*))`
	modelSplitRegexp          = `^\s*([A-Z][a-zA-Z]*)\s*{((?:.|\n)*?)}`
	attributeWithValueRegexp  = `^@([a-zA-Z]+)\(([^)]+)\)`
	attributeWithValueRegexp2 = `^@([a-zA-Z]+)\(([\s\S]+)\)`
)

type Model struct {
	SnakeName              string
	CamelName              string
	Columns                []*Column
	IsAuthRequired         bool
	IsAuthUser             bool
	ProtectedRoutes        []string
	HideRoutes             []string
	IsAuthOrganization     bool
	IsAuthOrganizationUser bool
	Hide                   bool
	Oauth2                 *Oauth2
}

type Oauth2 struct {
	Google bool
	Apple  bool
}

type ModelWithRelations struct {
	*Model
	Relations []*ModelTemplateRelation
}

type Column struct {
	SnakeName    string
	CamelName    string
	Type         *Type
	Attributes   []*Attribute
	IsRelation   bool
	Expose       bool
	DatabaseName *DatabaseName
}

type Type struct {
	Name                   string
	GoTypeName             string
	TypescriptName         string
	IsNullable             bool
	HasDifferentCharLength bool
	CharLength             int
	EmptyValue             string
}

type Attribute struct {
	Name     string
	Value    string
	HasValue bool
}

type TemplateData struct {
	PackageName string
}

type HelpersTemplateData struct {
	PackageName string
	JWTFields   []*JWTField
	HasAuth     bool
}

type JWTField struct {
	NormalName                  string
	CamelName                   string
	SnakeName                   string
	GoType                      string
	TableCamelName              string
	TableSnakeName              string
	IsFromUserTable             bool
	IsFromOrganizationTable     bool
	IsFromOrganizationUserTable bool
}

type TypescriptTypesTemplateData struct {
	Controllers []*Model
}

type TypescriptTypesRelationsTemplateData struct {
	Controllers           []*ModelWithRelations
	AuthFields            []*JWTField
	HasMultipleAuthFields bool
}

type SetupMainTemplateData struct {
	PackageName        string
	Imports            []string
	FullSetup          bool
	MigrationPath      string
	HasExtraMiddleWare bool
}

type GeneralTemplateData struct {
	PackageName string
	Controllers []*Model
	Imports     []string
}

type Oauth2TemplateData struct {
	PackageName           string
	UserTable             *Model
	OrganizationTable     *Model
	OrganizationUserTable *Model
	Imports               []string
	JWTFields             []*JWTField
}

type QueryTemplateData struct {
	PackageName           string
	Controllers           []*Model
	AuthFields            []*JWTField
	HasMultipleAuthFields bool
}

type RoutesTemplateData struct {
	PackageName           string
	Controllers           []*Model
	AuthFields            []*JWTField
	HasMultipleAuthFields bool
	HasOrganization       bool
}

type CreateAndUpdateDataModel struct {
	SnakeName     string
	CamelName     string
	CreateColumns []*Column
	UpdateColumns []*Column
	*Model
}

type CreateAndUpdateData struct {
	PackageName string
	Controllers []*CreateAndUpdateDataModel
	Imports     []string
	AuthField   *JWTField
}

type BodyTemplateData struct {
	PackageName string
	Controllers []*CreateAndUpdateDataModel
	Imports     []string
	AuthFields  []*JWTField
}

type SelectTemplateData struct {
	PackageName string
	Controllers []*ModelWithRelations
}

type ControllerTemplateData struct {
	PackageName   string
	CamelName     string
	Imports       []string
	CreateColumns []*Column
	UpdateColumns []*Column
	*Model
}

type AuthTemplateData struct {
	PackageName                   string
	CamelName                     string
	Imports                       []string
	CreateColumns                 []*Column
	JWTFields                     []*JWTField
	HasOrganization               bool
	OrganizationCamelName         string
	OrganizationCreateColumns     []*Column
	HasOrganizationUser           bool
	OrganizationUserCamelName     string
	OrganizationUserCreateColumns []*Column
}

type ModelTemplateRelation struct {
	Name         string
	SingularName string
	Columns      []*Column
	Type         string
	Tag          string
	IsArray      bool
	DatabaseName *DatabaseName
}

type DatabaseName struct {
	SingularCamelName string
	CamelName         string
}
type ModelTemplateData struct {
	PackageName string
	Imports     []string
	CamelName   string
	Relations   []*ModelTemplateRelation
	Columns     []*Column
}
