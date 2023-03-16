# SQL generator & REST api generator based on SQLboiler

Current implementation is only done for postgresql, though other sql databases like mysql should be easy to implement due to simular features. REST api is build on [gin-gonic](github.com/gin-gonic/gin).

## Installation

1. Create a new go project (`go mod new myproject`)
2. Create a new folder called `convert` in your project, add your `.env` file with `POSTGRESQL_URL` & add your `sqlboiler.toml` config file
3. Create a new file called `convert.go` in your `convert` folder
4. Copy the following code into your `convert.go` file:

```go
package main

import (
	"os"
	"strings"

	"github.com/arjendevos/gosql"
)

func main() {
	newDir, _ := os.Getwd()
	os.Chdir(strings.TrimSuffix(newDir, "/convert"))

	gosql.Convert(&gosql.GoSQLConfig{
		SchemeDir:           "schemes",
		MigrationDir:        "database/migrations",
		ModelOutputDir:      "models",
		ControllerOutputDir: "generated",
		SetupProject:        true,
	})
}
```

5. Run `go mod tidy` inside the `convert` folder
6. Create a new folder called `schemes` in your project
7. Create a new file called `1_migration.gosql` in your `schemes` folder (you can name it whatever you want, but the number is important)
8. Add your models into the `1_migration.gosql` file. Make sure to put `@postgresql` at the top (only needed for the first migration).
9. Run `(cd convert && go run convert.go)` in your project folder
10. Everything should be setup now, you can run `go run main.go` to start your server

## Possibilities

### Models

Example model:

```gosql
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

`@noAuth` is deprecated, use `@protected` instead. Use `@hide` to hide entire table from outside world.

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
- [x] filter={"column":{"equals":true, "or": {"lessThan": 4}}}
- [x] filter={"column":{"equals":true, "or": {"lessThan": 4, "isNotIn": ["d"]}}
- [x] filter={"relationTable": {"column": {"equals":true}}} (filter on relation)
- [x] rels={"relationTable":{}}
- [x] rels={"relationTable":{"deeperRelation":{ etc... }}} (you can't call the parent relation in the child relation)
- [x] rels={"relationTable":{"\_limit":4, "\_page":1}} (only for nToMany relations
- [x] order={"column": "desc", "column":"asc"}
- [x] from=organization | from=user | no parameter (organization = get by organization id, get = fetch by user id, no parameter = get by organization id & user id)
- [x] select=["column"] (omitempty fixes this on the json side)
- [x] select relation columns. {"\_select":["id", "created_at"]"}
- [ ] Add permissions to a role or user

### IDs

- [x] id uuid @unique @default(uuid_generate_v4()) @index
- [x] id int @unique @default(autoincrement) @index

## Todo

### Current implementations & future plans:

- [x] Add order by possibility
- [x] Add rest of the endpoints with bodies etc
- [x] Add or possibility in filter
- [x] Export to typescript types
- [x] Generate postgresql database setup files (client & migrations)
- [x] Auto install deps
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
- [x] limit queries to relations for x role
- [x] setup entire project
- [x] Test all endpoints (high priority)
- [x] middelware is somehow called 3 times (high priority)
- [ ] add enum for role (low priority)
- ~~ [ ] add filter to relation (low priority)~~
- [x] Add select columns on relation (\_select:["id", "name"])
- [x] Add select columns (normal not on relations)
- [x] Add oauth2 login option
- [x] Add oauth2 google login endpoints
- [x] Add custom marshaller for time.Time (high priority)
- [x] Make oauth2 google endpoints better based on scheme (org etc) (high priority)
- [x] Select columns fix created_at and updated_at omitempty problem (low priority)
- [ ] Add oauth2 facebook login endpoints (low priority)
- [ ] Add oauth2 apple login endpoints (lowest priority)
- [ ] add email option (smtp with default templates) (low priority)
- [ ] add password forget endpoints (very low priority) (should implement with email)

## Custom options

There is an option to add extra middleware in the auth middleware to handle role access. This will be generated automatically if you turn on `SetupProject`.
