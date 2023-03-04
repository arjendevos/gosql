# SQL generator & REST api generator based on SQLboiler

Current implementation is only done for postgresql, though other sql databases like mysql should be easy to implement due to simular features. REST api is build on [gin-gonic](github.com/gin-gonic/gin).

## Possibilities

### Models

Example model:

```
User {
  id uuid @unique @default(uuid_generate_v4()) @index
  name string
  email string @unique
  password string @hide
  organization Organization
  createdAt dateTime @default(now)
  updatedAt dateTime @default(now)
  deletedAt dateTime
}
```

### Authentication

Add `authUser` after your model to make it the main table for your user authentication. It will automatically add their relation to all tables. These columns are required:

- `email string @unique`
- `password string`

You can exclude any table from the authentication by adding `@noAuth` after the table name.

`@noAuth` is deprecated, use `@protected` instead.

`@protected` has the following options (`@protected(LIST, BYID, CREATE, UPDATE, DELETE)`):

- `LIST` - To protect the list endpoint
- `BYID` - To protect the by id endpoint
- `CREATE` - To protect the create endpoint
- `UPDATE` - To protect the update endpoint
- `DELETE` - To protect the delete endpoint

### Data types

| Type                    | PSQL Generator     | API Generator        |
| ----------------------- | ------------------ | -------------------- |
| `string`, `string(255)` | :white_check_mark: | :white_check_mark:   |
| `bool`                  | :white_check_mark: | :white_check_mark:   |
| `text`                  | :white_check_mark: | :white_check_mark:   |
| `dateTime`              | :white_check_mark: | :white_check_mark:   |
| `int`                   | :white_check_mark: | :white_check_mark:   |
| `any`                   | :white_check_mark: | :white_large_square: |

### Attributes

| Type                             | Meaning                                     |
| -------------------------------- | ------------------------------------------- |
| `?` after type                   | Is nullable                                 |
| `@uniue`                         | Is unique                                   |
| `@default(autoincrement)`        | Auto increment                              |
| `@default(uuid_generate_v4())`   | Auto Generate uuid                          |
| `@default(now)`                  | Auto generate current time                  |
| `@default("your default value")` | Default string value                        |
| `@default(false)`                | Default boolean value                       |
| `@default(1)`                    | Default int value                           |
| `@index`                         | Index on that column                        |
| `@hide`                          | Hide from outside world in the api          |
| `@regexp("your regexp")`         | Regexp validation for creating and updating |

### Relations

You can create a relation by adding a column name (this should be the table name with lowercase), like: `account` and you add as type the table name: `Account` with the first letter as capital. Don't refer both tables to each other, only one of them. Without `@unique` it is automatically a one to many relation. Optimize it by adding `@index` to the column.

- Add `?` to make it nullable.
- Add `@unique` to make it a one to one relation.

### Query parameters

- [x] limit=3
- [x] page=4
- [x] filter={"column": {"equals":"true"}}
- [x] rels={"relationTable":{}}
- [x] rels={"relationTable":{"deeperRelation":{ etc... }}}
- [x] order={"column": "desc", "column":"asc"}
- [x] filter={"column":{"equals":true, "or": {"lessThan": 4}}}
- [x] filter={"column":{"equals":true, "or": {"lessThan": 4, "isNotIn": ["d"]}}
- [x] filter={"relationTable": {"column": {"equals":true}}} // filter on relation
- [x] rels={"relationTable":{"\_limit":4, "\_page":1}} // limit relation array -> only for nToMany relations
- [x] from=organization | from=user | no parameter (organization = get by organization id, get = fetch by user id, no parameter = get by organization id & user id)
- [ ] select={"column":true"} // will become to much headache with golang to work (json fields)

### IDs

- [x] id uuid @unique @default(uuid_generate_v4()) @index
- [x] id int @unique @default(autoincrement) @index

## Todo

### Current implementations & future plans:

- [x] Add order by possibility
- [x] Add rest of the endpoints with bodies etc
- [ ] Test all endpoints
- [x] Add or possibility in filter
- [ ] Export to typescript types
- [ ] Generate postgresql database setup files (client & migrations)
- [ ] Auto install deps
- [x] Add authorization on User & Organization
- [x] User authentication
- [x] fetch items based on user_id or organization_id
- [x] expose relation ids in api
- [x] fetch relations for every request except create
- [x] change relations to include relations of relation
- [x] add pagination to relations
- [x] can filter on relation id's
- [x] can filter on if null
- [x] filter in relations in filter
- [x] fix if filter does not exists sql will output: WHERE ()
- [ ] Add select columns (normal not on relations) // will become to much headache with golang to work (json fields)

### Steps for setting up an api:

1. go get github.com/gin-gonic/gin
2. go get github.com/volatiletech/sqlboiler/v4
3. go get gopkg.in/validator.v2
4. go get github.com/volatiletech/null/v8
5. go get github.com/gin-contrib/cors
6. go get github.com/lib/pq
7. go get github.com/joho/godotenv
8. go get github.com/dgrijalva/jwt-go

## Custom options

There is an option to add extra middleware in the auth middleware to handle role access. For example:

!! There is 1 slight problem, because we can fetch relations by query paremeter, they could access things using relations. I'm not sure how to fix this yet.

```
var authCalls = map[string]map[string]bool{
	"admin": {
		"organization.byid":   true,
		"organization.list":   true,
		"organization.create": true,
		"organization.update": true,
		"organization.delete": true,
		"page.byid":           true,
		"page.list":           true,
	},
}

func CallMiddleware(c *gin.Context) {
	role := c.Value("role")

	method := strings.ToUpper(c.Request.Method)
	path := strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(c.Request.URL.Path), "/"), "/")
	pathSplitted := strings.Split(path, "/")

	var methodCall string

	switch method {
	case "GET":
		methodCall = "list"
		if contains(c.Params, "id") {
			methodCall = "byid"
			pathSplitted = pathSplitted[:len(pathSplitted)-1]
		}
	case "POST":
		methodCall = "create"
		if contains(c.Params, "id") {
			pathSplitted = pathSplitted[:len(pathSplitted)-1]
		}
	case "PATCH":
		methodCall = "update"
		if contains(c.Params, "id") {
			pathSplitted = pathSplitted[:len(pathSplitted)-1]
		}
	case "DELETE":
		methodCall = "delete"
		if contains(c.Params, "id") {
			pathSplitted = pathSplitted[:len(pathSplitted)-1]
		}
	}

	call := fmt.Sprintf("%s.%s", strings.Join(pathSplitted, "."), methodCall)

	if !authCalls[role.(string)][call] {
		fmt.Println("Call not allowed:")
		fmt.Println()
		fmt.Println("    ", call)
		fmt.Println()
		c.AbortWithStatusJSON(401, generated.ResponseWithPayload(nil, "unauthorized", "You are unauthorized", false))
		return
	}

	c.Next()
}

func contains(strSlice []gin.Param, str string) bool {
	for _, s := range strSlice {
		if s.Key == str {
			return true
		}
	}
	return false
}
```
