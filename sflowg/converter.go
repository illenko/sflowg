package sflowg

import (
	"fmt"
)

func ToStringValueMap(m map[string]any) map[string]string {
	result := make(map[string]string)
	for key, value := range m {
		switch v := value.(type) {
		case string:
			result[key] = v
		case int:
			result[key] = fmt.Sprintf("%d", v)
		case float64:
			result[key] = fmt.Sprintf("%f", v)
		case bool:
			result[key] = fmt.Sprintf("%t", v)
		case nil:
			result[key] = ""
		default:
			result[key] = fmt.Sprintf("%v", v)
		}
	}
	return result
}
