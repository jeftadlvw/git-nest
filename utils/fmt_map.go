package utils

import (
	"bytes"
	"fmt"
)

const defaultIndent string = "  "

func FmtMapTree(m map[string]interface{}, indent string) string {

	buffer := bytes.NewBufferString("")

	maxKeyLength := 0
	for key, _ := range m {
		maxKeyLength = max(maxKeyLength, len(key))
	}

	for key, value := range m {
		switch v := value.(type) {
		case map[string]interface{}:
			_, _ = fmt.Fprintf(buffer, "%s%s:\n", indent, key)
			formattedMapTree := FmtMapTree(v, indent+defaultIndent)
			_, _ = fmt.Fprintf(buffer, "%s\n", formattedMapTree)
		default:

			fmt.Fprintf(buffer, "%s%s:%-*s %v\n", indent, key, maxKeyLength-len(key)+2, "", value)
		}
	}

	return buffer.String()
}
