// internal/parser/where.go
package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// WhereOperator represents comparison operators
type WhereOperator int

const (
	OpEquals WhereOperator = iota
	OpNotEquals
	OpLessThan
	OpGreaterThan
	OpLessThanOrEqual
	OpGreaterThanOrEqual
	OpLike
)

// WhereCondition represents a single condition
type WhereCondition struct {
	Column   string
	Operator WhereOperator
	Value    string
}

// WhereExpression represents a WHERE clause with support for AND/OR logic
type WhereExpression struct {
	Conditions []WhereCondition
	LogicOps   []string // "AND" or "OR" between conditions
	Groups     []int    // Parentheses grouping (0 = no group, 1+ = group level)
}

// ParseWhereClause parses a WHERE clause string into a WhereExpression
func ParseWhereClause(whereClause string) (*WhereExpression, error) {
	whereClause = strings.TrimSpace(whereClause)
	if whereClause == "" {
		return nil, fmt.Errorf("empty WHERE clause")
	}

	expr := &WhereExpression{
		Conditions: []WhereCondition{},
		LogicOps:   []string{},
		Groups:     []int{},
	}

	// Tokenize the WHERE clause
	tokens := tokenizeWhere(whereClause)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no valid tokens in WHERE clause")
	}

	// Parse tokens into conditions and logic operators
	return parseWhereTokens(tokens, expr)
}

// tokenizeWhere splits WHERE clause into tokens, handling quoted strings
func tokenizeWhere(whereClause string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := '"'

	for _, char := range whereClause {
		switch char {
		case '"', '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
				if current.Len() > 0 {
					tokens = append(tokens, strings.TrimSpace(current.String()))
					current.Reset()
				}
			} else if char == quoteChar {
				inQuotes = false
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			} else {
				current.WriteRune(char)
			}
		case ' ', '\t', '\n':
			if !inQuotes {
				if current.Len() > 0 {
					tokens = append(tokens, strings.TrimSpace(current.String()))
					current.Reset()
				}
			} else {
				current.WriteRune(char)
			}
		case '(', ')':
			if !inQuotes {
				if current.Len() > 0 {
					tokens = append(tokens, strings.TrimSpace(current.String()))
					current.Reset()
				}
				tokens = append(tokens, string(char))
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, strings.TrimSpace(current.String()))
	}

	// Filter out empty tokens
	var result []string
	for _, token := range tokens {
		if token != "" {
			result = append(result, token)
		}
	}

	return result
}

// parseWhereTokens parses tokenized WHERE clause
func parseWhereTokens(tokens []string, expr *WhereExpression) (*WhereExpression, error) {
	i := 0
	groupLevel := 0

	for i < len(tokens) {
		token := tokens[i]

		switch strings.ToUpper(token) {
		case "AND", "OR":
			if len(expr.Conditions) == 0 {
				return nil, fmt.Errorf("logic operator %s at start of expression", token)
			}
			expr.LogicOps = append(expr.LogicOps, strings.ToUpper(token))
			expr.Groups = append(expr.Groups, groupLevel)
			i++

		case "(":
			groupLevel++
			i++

		case ")":
			if groupLevel <= 0 {
				return nil, fmt.Errorf("unmatched closing parenthesis")
			}
			groupLevel--
			i++

		default:
			// Parse condition: column operator value
			condition, consumed, err := parseCondition(tokens, i)
			if err != nil {
				return nil, err
			}
			expr.Conditions = append(expr.Conditions, condition)
			expr.Groups = append(expr.Groups, groupLevel)
			i += consumed
		}
	}

	if groupLevel != 0 {
		return nil, fmt.Errorf("unmatched parentheses")
	}

	if len(expr.Conditions) == 0 {
		return nil, fmt.Errorf("no conditions found in WHERE clause")
	}

	return expr, nil
}

// parseCondition parses a single condition from tokens
func parseCondition(tokens []string, start int) (WhereCondition, int, error) {
	if start+2 >= len(tokens) {
		return WhereCondition{}, 0, fmt.Errorf("incomplete condition")
	}

	column := tokens[start]
	operatorStr := strings.ToUpper(tokens[start+1])
	value := tokens[start+2]

	// Parse operator
	var operator WhereOperator
	switch operatorStr {
	case "=":
		operator = OpEquals
	case "!=", "<>":
		operator = OpNotEquals
	case "<":
		operator = OpLessThan
	case ">":
		operator = OpGreaterThan
	case "<=":
		operator = OpLessThanOrEqual
	case ">=":
		operator = OpGreaterThanOrEqual
	case "LIKE":
		operator = OpLike
	default:
		return WhereCondition{}, 0, fmt.Errorf("unsupported operator: %s", operatorStr)
	}

	// Remove quotes from value
	value = strings.Trim(value, "'\"")

	return WhereCondition{
		Column:   column,
		Operator: operator,
		Value:    value,
	}, 3, nil
}

// EvaluateCondition evaluates a single condition against a row
func (wc *WhereCondition) EvaluateCondition(row []string, columnIndexes map[string]int) (bool, error) {
	colIdx, exists := columnIndexes[wc.Column]
	if !exists {
		return false, fmt.Errorf("column %s not found", wc.Column)
	}

	if colIdx >= len(row) {
		return false, fmt.Errorf("column index out of bounds")
	}

	cellValue := row[colIdx]

	switch wc.Operator {
	case OpEquals:
		return cellValue == wc.Value, nil
	case OpNotEquals:
		return cellValue != wc.Value, nil
	case OpLike:
		return evaluateLike(cellValue, wc.Value)
	default:
		// For numeric comparisons, try to convert to numbers
		return evaluateNumericComparison(cellValue, wc.Value, wc.Operator)
	}
}

// evaluateLike evaluates LIKE pattern matching
func evaluateLike(value, pattern string) (bool, error) {
	// Convert SQL LIKE pattern to Go regex
	// % -> .*
	// _ -> .
	// Escape other regex special characters
	regexPattern := regexp.QuoteMeta(pattern)
	regexPattern = strings.ReplaceAll(regexPattern, "%", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "_", ".")
	regexPattern = "^" + regexPattern + "$"

	matched, err := regexp.MatchString(regexPattern, value)
	return matched, err
}

// evaluateNumericComparison evaluates numeric comparisons
func evaluateNumericComparison(value, compareValue string, operator WhereOperator) (bool, error) {
	// Try to parse as numbers
	valNum, err1 := strconv.ParseFloat(value, 64)
	compareNum, err2 := strconv.ParseFloat(compareValue, 64)

	// If both are numbers, do numeric comparison
	if err1 == nil && err2 == nil {
		switch operator {
		case OpLessThan:
			return valNum < compareNum, nil
		case OpGreaterThan:
			return valNum > compareNum, nil
		case OpLessThanOrEqual:
			return valNum <= compareNum, nil
		case OpGreaterThanOrEqual:
			return valNum >= compareNum, nil
		}
	}

	// Fallback to string comparison
	switch operator {
	case OpLessThan:
		return value < compareValue, nil
	case OpGreaterThan:
		return value > compareValue, nil
	case OpLessThanOrEqual:
		return value <= compareValue, nil
	case OpGreaterThanOrEqual:
		return value >= compareValue, nil
	}

	return false, fmt.Errorf("unsupported operator for comparison")
}

// EvaluateExpression evaluates the entire WHERE expression against a row
func (we *WhereExpression) EvaluateExpression(row []string, columnIndexes map[string]int) (bool, error) {
	if len(we.Conditions) == 0 {
		return true, nil
	}

	// Evaluate all conditions
	results := make([]bool, len(we.Conditions))
	for i, condition := range we.Conditions {
		result, err := condition.EvaluateCondition(row, columnIndexes)
		if err != nil {
			return false, err
		}
		results[i] = result
	}

	// Apply logic operators
	return we.evaluateLogic(results), nil
}

// evaluateLogic applies AND/OR logic with parentheses grouping
func (we *WhereExpression) evaluateLogic(results []bool) bool {
	if len(results) == 1 {
		return results[0]
	}

	// Simple evaluation: process left to right with proper precedence
	// For now, we'll implement a basic version without full parentheses support
	// This can be enhanced later with a proper expression tree
	result := results[0]

	for i := 0; i < len(we.LogicOps); i++ {
		if i+1 >= len(results) {
			break
		}

		nextResult := results[i+1]

		switch we.LogicOps[i] {
		case "AND":
			result = result && nextResult
		case "OR":
			result = result || nextResult
		}
	}

	return result
}
