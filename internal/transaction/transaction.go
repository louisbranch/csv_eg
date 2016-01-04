package transaction

import "time"

type Type int

const (
	Invalid Type = -1
	Debit   Type = iota
	Credit
)

type Transaction struct {
	Date        time.Time
	Description string
	Type
	Value uint
}
