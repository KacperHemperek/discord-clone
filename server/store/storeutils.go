package store

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	InvaliBoolFilterErr = errors.New("invalid bool filter string")
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
	return nil, InvaliBoolFilterErr
}
