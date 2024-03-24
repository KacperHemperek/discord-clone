package store

import (
	"database/sql"
	"fmt"
)

func rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		fmt.Println("Error occurred when rolling back: ", err)
	}
}
