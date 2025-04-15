package postgres

import "time"

type Operator string

const (
	OpEquals             Operator = "="
	OpNotEquals          Operator = "!="
	OpGreaterThan        Operator = ">"
	OpGreaterThanOrEqual Operator = ">="
	OpLessThan           Operator = "<"
	OpLessThanOrEqual    Operator = "<="
	OpBetween            Operator = "between"
	OpIn                 Operator = "in"
)

type TimeCondition map[Operator]time.Time

func FormatTimeCondition(sign Operator, date time.Time) TimeCondition {
	return TimeCondition{
		sign: date,
	}
}
