package models

import (
	"fmt"
	"strings"

	"github.com/onlysumitg/GoMockAPI/utils/typeutils"
)

var OperatorFuncMap = map[string]func(any, string, string) bool{
	"EQUALS_TO":                 EQUALS_TO,
	"NOT_EQUALS_TO":             NOT_EQUALS_TO,
	"LESS_THAN":                 LESS_THAN,
	"LESS_THAN_OR_EQUALS_TO":    LESS_THAN_OR_EQUALS_TO,
	"GREATER_THAN":              GREATER_THAN,
	"GREATER_THAN_OR_EQUALS_TO": GREATER_THAN_OR_EQUALS_TO,
	"CONTAINS":                  CONTAINS,
	"STARTS_WITH":               STARTS_WITH,
	"ENDS_WITH":                 ENDS_WITH,
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func EQUALS_TO(val1 any, val2 string, dataType string) bool {
	switch strings.ToUpper(dataType) {
	case "BOOL": // without quotes
		return typeutils.GetBoolVal(val1) == typeutils.GetBoolVal(val2)

	case "FLOAT64":
		return typeutils.GetFloatVal(val1) == typeutils.GetFloatVal(val2)

	case "INT":
		return typeutils.GetIntVal(val1) == typeutils.GetIntVal(val2)
	default:
		return strings.EqualFold(fmt.Sprint(val1), fmt.Sprint(val2))

	}

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func NOT_EQUALS_TO(val1 any, val2 string, dataType string) bool {
	return !EQUALS_TO(val1, val2, dataType)
}

// -----------------------------------------------------------------
// val1 < val2
// -----------------------------------------------------------------

func LESS_THAN(val1 any, val2 string, dataType string) bool {
	switch strings.ToUpper(dataType) {
	case "BOOL": // without quotes
		return typeutils.GetBoolVal(val1) != typeutils.GetBoolVal(val2)

	case "FLOAT64":
		return typeutils.GetFloatVal(val1) < typeutils.GetFloatVal(val2)
	case "INT":
		return typeutils.GetIntVal(val1) < typeutils.GetIntVal(val2)
	default:
		return strings.ToUpper(fmt.Sprint(val1)) < strings.ToUpper(fmt.Sprint(val2))
	}
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------

func LESS_THAN_OR_EQUALS_TO(val1 any, val2 string, dataType string) bool {
	switch strings.ToUpper(dataType) {
	case "BOOL": // without quotes
		return typeutils.GetBoolVal(val1) == typeutils.GetBoolVal(val2)

	case "FLOAT64":
		return typeutils.GetFloatVal(val1) <= typeutils.GetFloatVal(val2)
	case "INT":
		return typeutils.GetIntVal(val1) <= typeutils.GetIntVal(val2)
	default:
		return strings.ToUpper(fmt.Sprint(val1)) <= strings.ToUpper(fmt.Sprint(val2))
	}

}

// -----------------------------------------------------------------
// val1 > val2
// -----------------------------------------------------------------

func GREATER_THAN(val1 any, val2 string, dataType string) bool {
	switch strings.ToUpper(dataType) {
	case "BOOL": // without quotes
		return typeutils.GetBoolVal(val1) != typeutils.GetBoolVal(val2)

	case "FLOAT64":
		return typeutils.GetFloatVal(val1) > typeutils.GetFloatVal(val2)
	case "INT":
		return typeutils.GetIntVal(val1) > typeutils.GetIntVal(val2)
	default:
		return strings.ToUpper(fmt.Sprint(val1)) > strings.ToUpper(fmt.Sprint(val2))
	}
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------

func GREATER_THAN_OR_EQUALS_TO(val1 any, val2 string, dataType string) bool {
	switch strings.ToUpper(dataType) {
	case "BOOL": // without quotes
		return typeutils.GetBoolVal(val1) == typeutils.GetBoolVal(val2)

	case "FLOAT64":
		return typeutils.GetFloatVal(val1) >= typeutils.GetFloatVal(val2)
	case "INT":
		return typeutils.GetIntVal(val1) >= typeutils.GetIntVal(val2)
	default:
		return strings.ToUpper(fmt.Sprint(val1)) >= strings.ToUpper(fmt.Sprint(val2))
	}

}

// -----------------------------------------------------------------
// val1 contains val2
// -----------------------------------------------------------------

func CONTAINS(val1 any, val2 string, dataType string) bool {

	return strings.Contains(strings.ToUpper(fmt.Sprint(val1)), strings.ToUpper(fmt.Sprint(val2)))

}

// -----------------------------------------------------------------
// val1 starts with val2
// -----------------------------------------------------------------

func STARTS_WITH(val1 any, val2 string, dataType string) bool {
	return strings.HasPrefix(strings.ToUpper(fmt.Sprint(val1)), strings.ToUpper(fmt.Sprint(val2)))

}

// -----------------------------------------------------------------
// val1 ends with val2
// -----------------------------------------------------------------

func ENDS_WITH(val1 any, val2 string, dataType string) bool {
	return strings.HasSuffix(strings.ToUpper(fmt.Sprint(val1)), strings.ToUpper(fmt.Sprint(val2)))

}
