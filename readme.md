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

### Current query possibility:

limit=3
page=4
filter={"hello": {"equals":"true"}}
rels=["job": {"id": true, "name":true, "helloWorld":true}]

### Todo:

select={"id":true}
filter={"hello":{"equals":true, "or": {"lessThan": 4}}}
filter={"hello":{"equals":true, "and": {"lessThan": 4}}}
filter={"hello":{"equals":true, "or": {"lessThan": 4, "and": {"isNotIn": ["d"]}}}}
order=["name": "desc", "id":"ac"]
