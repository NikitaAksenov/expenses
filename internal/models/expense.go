package models

import (
	"fmt"
	"time"

	"github.com/NikitaAksenov/expenses/internal"
)

type Expense struct {
	ID       int
	Name     string
	Category string
	Value    float64
	Date     time.Time
	Created  time.Time
	Updated  time.Time
}

func (e Expense) String() string {
	return fmt.Sprintf("#%03d | %15s | %15s | %7.2f | %s | %s | %s", e.ID, e.Name, e.Category, e.Value, e.Date.Format(internal.TimeFormatDateOnly), e.Created.Format(internal.TimeFormatDateOnly), e.Updated.Format(internal.TimeFormatDateOnly))
}
