# Model types

## IDs

1. id uuid @unique @default(uuid_generate_v4()) @index
2. id int @unique @default(autoincrement) @index

## Types

1. string
2. bool
3. text
4. DateTime
5. int

## TODO

1. Add order by possibility
2. Add select columns (normal not on relations)
3. Add rest of the endpoints with bodies etc
4. Add or & and possibility in filter
5. Export to typescript types

### Current query possibility:

- limit=3
- page=4
- filter={"column": {"equals":"true"}}
- rels={"relationTable": {"column": true, "column2":true, "column3":true}}
- rels={"relationTable":{}} // fetch all
- rels={"relationTable":{"_all":true}} // fetch all

### Todo:

- select={"column":true}
- filter={"column":{"equals":true, "or": {"lessThan": 4}}}
- filter={"column":{"equals":true, "and": {"lessThan": 4}}}
- filter={"column":{"equals":true, "or": {"lessThan": 4, "and": {"isNotIn": ["d"]}}}}
- order=["column": "desc", "column":"ac"]
- rels={"relationTable":{"_limit":4, "_page":1}} // limit relation array

### Steps for suc6:

1. go get github.com/gin-gonic/gin
2. go get github.com/volatiletech/sqlboiler/v4
3. go get gopkg.in/validator.v2
4. go get github.com/volatiletech/null/v8
5. go get github.com/gin-contrib/cors

6. go get github.com/lib/pq
7. go get github.com/joho/godotenv
