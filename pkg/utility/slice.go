package utility

import (
	"fmt"
	"strings"
)

func FloatSliceToCommaSlice[T float32 | float64](ids []T) string {
	var sb strings.Builder

	for i, id := range ids {
		sb.WriteString(fmt.Sprintf("%f", id))
		if i < len(ids)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}
