package userup

type LogicalOperator string

const (
	And LogicalOperator = "$and"
	Or  LogicalOperator = "$or"
)

type Join struct {
	Table  string               `json:"table"`
	On     string               `json:"on"`
	Filter map[string]Condition `json:"filter,omitempty"`
}

type Condition map[string]interface{}

type Query struct {
	Filter  map[string]Condition `json:"filter"`
	Select  []string             `json:"select,omitempty"`
	OrderBy []Order              `json:"order_by,omitempty"`
	Limit   int                  `json:"limit,omitempty"`
	Offset  int                  `json:"offset,omitempty"`
	Joins   []Join               `json:"joins,omitempty"`
}

// Order specifies the ordering of the query results.
type Order struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // "ASC" or "DESC"
}

