package userup

// LogicalOperator represents a logical operator used in query conditions.
type LogicalOperator string

// Join represents a join operation in a database query.
type Join struct {
	Table  string               `json:"table"`            // The name of the table to join.
	On     string               `json:"on"`               // The join condition.
	Filter map[string]Condition `json:"filter,omitempty"` // Optional filter conditions for the join.
}

// Condition represents a condition in a database query.
type Condition map[string]interface{}

// Query represents a database query.
type Query struct {
	Filter  map[string]Condition `json:"filter"`             // The filter conditions for the query.
	Select  []string             `json:"select,omitempty"`   // The fields to select in the query.
	OrderBy []Order              `json:"order_by,omitempty"` // The ordering of the query results.
	Limit   int                  `json:"limit,omitempty"`    // The maximum number of results to return.
	Offset  int                  `json:"offset,omitempty"`   // The offset of the query results.
	Joins   []Join               `json:"joins,omitempty"`    // The join operations in the query.
}

// Order specifies the ordering of the query results.
type Order struct {
	Field     string `json:"field"`     // The field to order by.
	Direction string `json:"direction"` // The direction of the ordering ("ASC" or "DESC").
}
