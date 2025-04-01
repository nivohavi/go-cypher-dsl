package expr

import (
	"fmt"
	"strconv"
	"strings"
)

// formatValue formats a value to be used in a Cypher query
func formatValue(value any) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "\\'") + "'"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return strconv.FormatBool(v)
	case []any:
		parts := make([]string, len(v))
		for i, item := range v {
			parts[i] = formatValue(item)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]any:
		parts := make([]string, 0, len(v))
		for key, val := range v {
			parts = append(parts, key+": "+formatValue(val))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		return fmt.Sprintf("%v", v)
	}
}
