package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

var (
	InvalidBoolFilterErr   = errors.New("invalid bool filter string")
	InvalidNumberFilterErr = errors.New("invalid number filter string")
	LimitNumberTooSmallErr = errors.New("limit number must be at least 1")
)

func rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		fmt.Println("Error occurred when rolling back: ", err)
	}
}

func whereSQL(list []string) string {
	result := ""
	if len(list) == 0 {
		return result
	}
	for i, clause := range list {
		if i != 0 {
			result = result + " AND " + clause
			continue
		}
		result = "WHERE " + clause
	}
	return result
}

type BoolFilter = string

func NewBoolFilter(val string) (*BoolFilter, error) {
	if val == "true" || val == "false" {
		return &val, nil
	}
	return nil, InvalidBoolFilterErr
}

type LimitFilter = int

func NewLimitFilter(val any) (*LimitFilter, error) {
	switch v := val.(type) {
	case int:
		if v < 1 {
			return &v, LimitNumberTooSmallErr
		}
		return &v, nil
	case string:
		if v == "" {
			return nil, nil
		}
		intVal, err := strconv.Atoi(v)
		if intVal < 0 {
			return &intVal, LimitNumberTooSmallErr
		}
		return &intVal, err
	default:
		return nil, InvalidNumberFilterErr
	}

}
