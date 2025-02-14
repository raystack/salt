package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/raystack/salt/qp/parser"
)

type Organization struct {
	Id              int       `qp:"type=number,min=10,max=200"`
	BillingPlanName string    `qp:"type=string"`
	CreatedAt       time.Time `qp:"type=datetime"`
	MemberCount     int       `qp:"type=number"`
	Title           string    `qp:"type=string"`
	Enabled         bool      `qp:"type=bool"`
}

func main() {
	userInput := &parser.Query{}
	// org := Organization{10, "standard plan", time.Now(), 10, "pixxel space pvt ltd"}
	jsonFile, err := os.Open("input.json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = json.Unmarshal(byteValue, userInput)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal query string to parser query struct, err:%s", err.Error()))
	}

	err = parser.ValidateQuery(userInput, Organization{})
	if err != nil {
		panic(err)
	}
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

}
