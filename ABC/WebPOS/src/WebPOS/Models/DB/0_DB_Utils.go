package Models

import (
	"strings"
)

// Trim all character '%' : prefix & suffix
func trimCondition(condition string) string {
	trimCondition := strings.TrimSpace(condition)
	trimCondition = strings.TrimPrefix(trimCondition, _LIKE_CHAR)
	trimCondition = strings.TrimSuffix(trimCondition, _LIKE_CHAR)
	return trimCondition
}

// Add prefix & suffix character '%'
func fullLike(condition string) string {
	condition = trimCondition(condition)
	return "%" + condition + "%"
}

// Add prefix character '%'
func prefixLike(condition string) string {
	condition = trimCondition(condition)
	return "%" + condition
}

// Add suffix character '%'
func suffixLike(condition string) string {
	condition = trimCondition(condition)
	return condition + "%"
}
