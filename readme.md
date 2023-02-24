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
- `password` string

You can exclude any table from the authentication by adding `@noAuth` after the table name.

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

You can create a relation by adding a column name it whatever you want, like: `account` and you add as type the table name: `Account` with the first letter as capital. Don't refer both tables to each other, only one of them. Without `@unique` it is automatically a one to many relation. Optimize it by adding `@index` to the column.

- Add `?` to make it nullable.
- Add `@unique` to make it a one to one relation.

### IDs

- [x] id uuid @unique @default(uuid_generate_v4()) @index
- [x] id int @unique @default(autoincrement) @index

### TODO

- [x] Add order by possibility
- [x] Add rest of the endpoints with bodies etc
- [ ] Add or & and possibility in filter
- [ ] Export to typescript types
- [ ] Generate postgresql database setup files (client & migrations)
- [ ] Auto install deps
- [ ] Add authorization on User & Organization
- [ ] User authentication
- [ ] Add select columns (normal not on relations) // Not possible for now due to the sqlboiler implementation

### REST api -> query possibilities:

- [x] limit=3
- [x] page=4
- [x] filter={"column": {"equals":"true"}}
- [x] rels={"relationTable":true}
- [x] order={"column": "desc", "column":"asc"}
- [ ] filter={"column":{"equals":true, "or": {"lessThan": 4}}}
- [ ] filter={"column":{"equals":true, "and": {"lessThan": 4}}}
- [ ] filter={"column":{"equals":true, "or": {"lessThan": 4, "and": {"isNotIn": ["d"]}}}}
- [ ] rels={"relationTable":{"\_limit":4, "\_page":1}} // limit relation array -> only for nToMany relations
- [ ] select={"column":true"}

### Steps for setting up an api:

1. go get github.com/gin-gonic/gin
2. go get github.com/volatiletech/sqlboiler/v4
3. go get gopkg.in/validator.v2
4. go get github.com/volatiletech/null/v8
5. go get github.com/gin-contrib/cors
6. go get github.com/lib/pq
7. go get github.com/joho/godotenv
8. go get github.com/dgrijalva/jwt-go
