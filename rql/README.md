# RQL (Rest Query Language)

A library to parse support advanced REST API query parameters like (filter, pagination, sort, group, search) and logical operators on the keys (like eq, neq, like, gt, lt etc)

It takes a Golang struct and a json string as input and returns a Golang object that can be used to prepare SQL Statements (using raw sql or ORM Query builders).

### Usage

Frontend should send the parameters and operator like this schema to the backend service on some route with `POST` HTTP Method

```json
{
  "filters": [
    { "name": "id", "operator": "neq", "value": 20 },
    { "name": "title", "operator": "neq", "value": "nasa" },
    { "name": "enabled", "operator": "eq", "value": false },
    {
      "name": "created_at",
      "operator": "gte",
      "value": "2025-02-05T11:25:37.957Z"
    },
    { "name": "title", "operator": "like", "value": "xyz" }
  ],
  "group_by": ["billing_plan_name"],
  "offset": 20,
  "limit": 50,
  "search": "abcd",
  "sort": [
    { "key": "title", "order": "desc" },
    { "key": "created_at", "order": "asc" }
  ]
}
```

The `rql` library can be used to parse this json, validate it and returns a Struct containing all the info to generate the operations and values for SQL.

The validation happens via stuct tags defined on your model. Example:

```golang
type Organization struct {
	Id              int       `rql:"type=number,min=10,max=200"`
	BillingPlanName string    `rql:"type=string"`
	CreatedAt       time.Time `rql:"type=datetime"`
	MemberCount     int       `rql:"type=number"`
	Title           string    `rql:"type=string"`
	Enabled         bool      `rql:"type=bool"`
}

```

**Supported data types:**

1. number
2. string
3. datetime
4. bool

Check `main.go` for more info on usage.

Using this struct, a SQL query can be generated. Here is an example using `goqu` SQL Builder

```go
	//init the library's "Query" object with input json bytes
	userInput := &parser.Query{}

	//assuming jsonBytes is defined earlier
	err = json.Unmarshal(jsonBytes, userInput)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal query string to parser query struct, err:%s", err.Error()))
	}

	//validate the json input
	err = parser.ValidateQuery(userInput, Organization{})
	if err != nil {
		panic(err)
	}

	//userInput object can be utilized to prepare SQL statement
	query := goqu.From("organizations")

	fuzzySearchColumns := []string{"id", "billing_plan_name", "title"}

	for _, filter_item := range userInput.Filters {
		query = query.Where(goqu.Ex{
			filter_item.Name: goqu.Op{filter_item.Operator: filter_item.Value},
		})
	}

	listOfExpressions := make([]goqu.Expression, 0)

	if userInput.Search != "" {
		for _, col := range fuzzySearchColumns {
			listOfExpressions = append(listOfExpressions, goqu.Ex{
				col: goqu.Op{"LIKE": userInput.Search},
			})
		}
	}

	query = query.Where(goqu.Or(listOfExpressions...))

	query = query.Offset(uint(userInput.Offset))
	for _, sort_item := range userInput.Sort {
		switch sort_item.Order {
		case "asc":
			query = query.OrderAppend(goqu.C(sort_item.Key).Asc())
		case "desc":
			query = query.OrderAppend(goqu.C(sort_item.Key).Desc())
		default:
		}
	}
	query = query.Limit(uint(userInput.Limit))
	sql, _, _ := query.ToSQL()
	fmt.Println(sql)


```

giving output as

```sql
SELECT * FROM "organizations" WHERE (("id" != 20) AND ("title" != 'nasa') AND ("enabled" IS FALSE) AND ("createdAt" >= '2025-02-05T11:25:37.957Z') AND ("title" LIKE 'xyz') AND (("id" LIKE 'abcd') OR ("billing_plan_name" LIKE 'abcd') OR ("title" LIKE 'abcd'))) ORDER BY "title" DESC, "createdAt" ASC LIMIT 50 OFFSET 20
```

### Improvements

1. The operators need to mapped with SQL operators like (`eq` should be converted to `=` etc). Right now we are relying on GoQU to do that, but we can make it SQL ORL lib agnostic.

2. Support validation on the range or values of the data. Like `min`, `max` on number etc.
