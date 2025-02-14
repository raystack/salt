package parser

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/raystack/salt/utils"
)

var validNumberOperations = []string{"eq", "neq", "gt", "gte", "lte"}
var validStringOperations = []string{"eq", "neq", "like"}
var validBoolOperations = []string{"eq", "neq"}
var validDatetimeOperations = []string{"eq", "neq", "gt", "gt", "gte", "lte"}
var validSortOrder = []string{"asc", "desc"}

const TAG = "qp"

type Query struct {
	Filters []Filter `json:"filters"`
	GroupBy []string `json:"group_by"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Search  string   `json:"search"`
	Sort    []Sort   `json:"sort"`
}

type Filter struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Type     string `json:"type"`
	Value    any    `json:"value"`
}

type Sort struct {
	Key   string `json:"key"`
	Order string `json:"order"`
}

func ValidateQuery(q *Query, checkStruct interface{}) error {
	val := reflect.ValueOf(checkStruct)

	// validate filters
	for _, filterItem := range q.Filters {
		//validate filter key name
		filterIdx := searchKeyInsideStruct(filterItem.Name, val)
		if filterIdx < 0 {
			return fmt.Errorf("'%s' is not a valid filter key", filterItem.Name)
		}
		structKeyTag := val.Type().Field(filterIdx).Tag.Get(TAG)

		// validate filter key data type
		allowedDataType := getDataTypeOfField(structKeyTag)
		filterItem.Type = allowedDataType
		switch allowedDataType {
		case "number":
			err := validateNumberType(filterItem)
			if err != nil {
				return err
			}
		case "bool":
			err := validateBoolType(filterItem)
			if err != nil {
				return err
			}
		case "datetime":
			err := validateDatetimeType(filterItem)
			if err != nil {
				return err
			}
		case "string":
			err := validateStringType(filterItem)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("type '%s' is not recognized", allowedDataType)
		}

		if !isValidOperator(filterItem) {
			return fmt.Errorf("value '%s' for key '%s' is valid string", filterItem.Operator, filterItem.Name)
		}
	}

	err := validateGroupByKeys(q, val)
	if err != nil {
		return err
	}
	return validateSortKey(q, val)
}

func validateNumberType(filterItem Filter) error {
	// check if the type is any of Golang numeric types
	// if not, return error
	switch filterItem.Value.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64, int, uint:
		return nil
	default:
		return fmt.Errorf("value %v for key '%s' is not int type", filterItem.Value, filterItem.Name)
	}
}

func validateDatetimeType(filterItem Filter) error {
	// cast the value to datetime
	// if failed, return error
	castedVal, ok := filterItem.Value.(string)
	if !ok {
		return fmt.Errorf("value %s for key '%s' is not a valid ISO datetime string", filterItem.Value, filterItem.Name)
	}
	_, err := time.Parse(time.RFC3339, castedVal)
	if err != nil {
		return fmt.Errorf("value %s for key '%s' is not a valid ISO datetime string", filterItem.Value, filterItem.Name)
	}
	return nil
}

func validateBoolType(filterItem Filter) error {
	// cast the value to bool
	// if failed, return error
	_, ok := filterItem.Value.(bool)
	if !ok {
		return fmt.Errorf("value %v for key '%s' is not bool type", filterItem.Value, filterItem.Name)
	}
	return nil
}

func validateStringType(filterItem Filter) error {
	// cast the value to string
	// if failed, return error
	_, ok := filterItem.Value.(string)
	if !ok {
		return fmt.Errorf("value %s for key '%s' is valid string type", filterItem.Value, filterItem.Name)
	}
	return nil
}

func searchKeyInsideStruct(keyName string, val reflect.Value) int {
	for i := 0; i < val.NumField(); i++ {
		if strings.ToLower(val.Type().Field(i).Name) == strings.ToLower(keyName) {
			return i
		}
	}
	return -1
}

// parse the tag schema which is of the format
// type=int,min=10,max=200
// to extract type else fallback to string
func getDataTypeOfField(tagString string) string {
	res := "string"
	splitted := strings.Split(tagString, ",")
	for _, item := range splitted {
		kvSplitted := strings.Split(item, "=")
		if len(kvSplitted) == 2 {
			if kvSplitted[0] == "type" {
				return kvSplitted[1]
			}
		}
	}
	//fallback to string if type not found in tag value
	return res
}

func isValidOperator(filterItem Filter) bool {
	switch filterItem.Type {
	case "number":
		return utils.StringFoundInArray(filterItem.Operator, validNumberOperations)
	case "datetime":
		return utils.StringFoundInArray(filterItem.Operator, validDatetimeOperations)
	case "string":
		return utils.StringFoundInArray(filterItem.Operator, validStringOperations)
	case "bool":
		return utils.StringFoundInArray(filterItem.Operator, validBoolOperations)
	default:
		return false
	}
}

func validateSortKey(q *Query, val reflect.Value) error {
	for _, item := range q.Sort {
		filterIdx := searchKeyInsideStruct(item.Key, val)
		if filterIdx < 0 {
			return fmt.Errorf("'%s' is not a valid sort key", item.Key)
		}
		if !utils.StringFoundInArray(item.Order, validSortOrder) {
			return fmt.Errorf("'%s' is not a valid sort key", item.Key)
		}
	}
	return nil
}

func validateGroupByKeys(q *Query, val reflect.Value) error {
	for _, item := range q.GroupBy {
		filterIdx := searchKeyInsideStruct(item, val)
		if filterIdx < 0 {
			return fmt.Errorf("'%s' is not a valid sort key", item)
		}
	}
	return nil
}
