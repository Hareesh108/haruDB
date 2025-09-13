// internal/parser/where_test.go
package parser

import (
	"testing"
)

func TestParseWhereClause(t *testing.T) {
	tests := []struct {
		name        string
		whereClause string
		expectError bool
		expectedOps []WhereOperator
	}{
		{
			name:        "simple equality",
			whereClause: "age = 25",
			expectError: false,
			expectedOps: []WhereOperator{OpEquals},
		},
		{
			name:        "not equals",
			whereClause: "name != 'John'",
			expectError: false,
			expectedOps: []WhereOperator{OpNotEquals},
		},
		{
			name:        "less than",
			whereClause: "age < 30",
			expectError: false,
			expectedOps: []WhereOperator{OpLessThan},
		},
		{
			name:        "greater than",
			whereClause: "salary > 50000",
			expectError: false,
			expectedOps: []WhereOperator{OpGreaterThan},
		},
		{
			name:        "less than or equal",
			whereClause: "age <= 25",
			expectError: false,
			expectedOps: []WhereOperator{OpLessThanOrEqual},
		},
		{
			name:        "greater than or equal",
			whereClause: "score >= 80",
			expectError: false,
			expectedOps: []WhereOperator{OpGreaterThanOrEqual},
		},
		{
			name:        "like pattern",
			whereClause: "name LIKE 'John%'",
			expectError: false,
			expectedOps: []WhereOperator{OpLike},
		},
		{
			name:        "and condition",
			whereClause: "age > 18 AND status = 'active'",
			expectError: false,
			expectedOps: []WhereOperator{OpGreaterThan, OpEquals},
		},
		{
			name:        "or condition",
			whereClause: "age < 18 OR age > 65",
			expectError: false,
			expectedOps: []WhereOperator{OpLessThan, OpGreaterThan},
		},
		{
			name:        "complex and or",
			whereClause: "age > 18 AND (status = 'active' OR role = 'admin')",
			expectError: false,
			expectedOps: []WhereOperator{OpGreaterThan, OpEquals, OpEquals},
		},
		{
			name:        "empty clause",
			whereClause: "",
			expectError: true,
		},
		{
			name:        "incomplete condition",
			whereClause: "age =",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseWhereClause(tt.whereClause)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(expr.Conditions) != len(tt.expectedOps) {
				t.Errorf("expected %d conditions, got %d", len(tt.expectedOps), len(expr.Conditions))
				return
			}
			for i, expectedOp := range tt.expectedOps {
				if expr.Conditions[i].Operator != expectedOp {
					t.Errorf("condition %d: expected operator %v, got %v", i, expectedOp, expr.Conditions[i].Operator)
				}
			}
		})
	}
}

func TestEvaluateCondition(t *testing.T) {
	columnIndexes := map[string]int{
		"name":  0,
		"age":   1,
		"email": 2,
	}

	tests := []struct {
		name      string
		row       []string
		condition WhereCondition
		expected  bool
	}{
		{
			name:      "equals match",
			row:       []string{"John", "25", "john@example.com"},
			condition: WhereCondition{Column: "name", Operator: OpEquals, Value: "John"},
			expected:  true,
		},
		{
			name:      "equals no match",
			row:       []string{"John", "25", "john@example.com"},
			condition: WhereCondition{Column: "name", Operator: OpEquals, Value: "Jane"},
			expected:  false,
		},
		{
			name:      "not equals",
			row:       []string{"John", "25", "john@example.com"},
			condition: WhereCondition{Column: "name", Operator: OpNotEquals, Value: "Jane"},
			expected:  true,
		},
		{
			name:      "numeric less than",
			row:       []string{"John", "25", "john@example.com"},
			condition: WhereCondition{Column: "age", Operator: OpLessThan, Value: "30"},
			expected:  true,
		},
		{
			name:      "numeric greater than",
			row:       []string{"John", "25", "john@example.com"},
			condition: WhereCondition{Column: "age", Operator: OpGreaterThan, Value: "20"},
			expected:  true,
		},
		{
			name:      "like pattern match",
			row:       []string{"John", "25", "john@example.com"},
			condition: WhereCondition{Column: "email", Operator: OpLike, Value: "john%"},
			expected:  true,
		},
		{
			name:      "like pattern no match",
			row:       []string{"John", "25", "john@example.com"},
			condition: WhereCondition{Column: "email", Operator: OpLike, Value: "jane%"},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.condition.EvaluateCondition(tt.row, columnIndexes)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluateExpression(t *testing.T) {
	columnIndexes := map[string]int{
		"name":   0,
		"age":    1,
		"status": 2,
	}

	tests := []struct {
		name     string
		row      []string
		expr     *WhereExpression
		expected bool
	}{
		{
			name: "single condition true",
			row:  []string{"John", "25", "active"},
			expr: &WhereExpression{
				Conditions: []WhereCondition{
					{Column: "age", Operator: OpGreaterThan, Value: "20"},
				},
				LogicOps: []string{},
				Groups:   []int{0},
			},
			expected: true,
		},
		{
			name: "single condition false",
			row:  []string{"John", "25", "active"},
			expr: &WhereExpression{
				Conditions: []WhereCondition{
					{Column: "age", Operator: OpLessThan, Value: "20"},
				},
				LogicOps: []string{},
				Groups:   []int{0},
			},
			expected: false,
		},
		{
			name: "and condition both true",
			row:  []string{"John", "25", "active"},
			expr: &WhereExpression{
				Conditions: []WhereCondition{
					{Column: "age", Operator: OpGreaterThan, Value: "20"},
					{Column: "status", Operator: OpEquals, Value: "active"},
				},
				LogicOps: []string{"AND"},
				Groups:   []int{0, 0},
			},
			expected: true,
		},
		{
			name: "and condition one false",
			row:  []string{"John", "25", "inactive"},
			expr: &WhereExpression{
				Conditions: []WhereCondition{
					{Column: "age", Operator: OpGreaterThan, Value: "20"},
					{Column: "status", Operator: OpEquals, Value: "active"},
				},
				LogicOps: []string{"AND"},
				Groups:   []int{0, 0},
			},
			expected: false,
		},
		{
			name: "or condition one true",
			row:  []string{"John", "25", "inactive"},
			expr: &WhereExpression{
				Conditions: []WhereCondition{
					{Column: "age", Operator: OpGreaterThan, Value: "20"},
					{Column: "status", Operator: OpEquals, Value: "active"},
				},
				LogicOps: []string{"OR"},
				Groups:   []int{0, 0},
			},
			expected: true,
		},
		{
			name: "or condition both false",
			row:  []string{"John", "15", "inactive"},
			expr: &WhereExpression{
				Conditions: []WhereCondition{
					{Column: "age", Operator: OpGreaterThan, Value: "20"},
					{Column: "status", Operator: OpEquals, Value: "active"},
				},
				LogicOps: []string{"OR"},
				Groups:   []int{0, 0},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.expr.EvaluateExpression(tt.row, columnIndexes)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTokenizeWhere(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple condition",
			input:    "age = 25",
			expected: []string{"age", "=", "25"},
		},
		{
			name:     "quoted string",
			input:    "name = 'John Doe'",
			expected: []string{"name", "=", "John Doe"},
		},
		{
			name:     "double quoted string",
			input:    `name = "John Doe"`,
			expected: []string{"name", "=", "John Doe"},
		},
		{
			name:     "and condition",
			input:    "age > 18 AND status = 'active'",
			expected: []string{"age", ">", "18", "AND", "status", "=", "active"},
		},
		{
			name:     "parentheses",
			input:    "(age > 18) AND (status = 'active')",
			expected: []string{"(", "age", ">", "18", ")", "AND", "(", "status", "=", "active", ")"},
		},
		{
			name:     "like with spaces",
			input:    "name LIKE 'John%' AND age > 25",
			expected: []string{"name", "LIKE", "John%", "AND", "age", ">", "25"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizeWhere(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("token %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}
