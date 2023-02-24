# SQL generator & REST api generator based on SQLboiler

Current implementation is only done for postgresql, though other sql databases like mysql should be easy to implement due to simular features. REST api is build on [gin-gonic](github.com/gin-gonic/gin).

## Model types

| Type                    | PSQL Generator | API Generator |
| ----------------------- | -------------- | ------------- |
| `string`, `string(255)` | - [x]          | - [x]         |
| `bool`                  | - [x]          | - [x]         |
| `text`                  | - [x]          | - [x]         |
| `dateTime`              | - [x]          | - [x]         |
| `int`                   | - [x]          | - [x]         |
| `any`                   | - [x]          | - [ ]         |

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
