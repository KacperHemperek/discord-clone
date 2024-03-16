package models

import (
	"database/sql"
	"fmt"
	"time"
)

type NullTime sql.NullTime

func (t *NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	date := fmt.Sprintf("\"%s\"", t.Time.Format(time.RFC3339))

	return []byte(date), nil
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	dataStr := string(data)

	if dataStr == "null" {
		t.Valid = false
		return nil
	}

	t.Valid = true
	timeFromString, err := time.Parse(time.RFC3339, dataStr)

	if err != nil {
		return err
	}

	t.Time = timeFromString
	return nil
}

func (t *NullTime) Scan(value any) error {
	nullTime := sql.NullTime{}

	err := nullTime.Scan(value)

	if err != nil {
		return err
	}

	*t = NullTime(nullTime)

	return nil
}
