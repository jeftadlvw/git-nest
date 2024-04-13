package utils

import (
	"bytes"
	"fmt"
)

type Node struct {
	Key   string
	Value interface{}
}

const defaultIndent string = "  "

func FmtTree(indent string, nodes ...Node) string {

	buffer := bytes.NewBufferString("")

	maxKeyLength := 0
	for _, node := range nodes {
		maxKeyLength = max(maxKeyLength, len(node.Key))
	}

	for _, node := range nodes {
		switch v := node.Value.(type) {
		case []Node:
			_, _ = fmt.Fprintf(buffer, "%s%s:\n", indent, node.Key)
			formattedMapTree := FmtTree(indent+defaultIndent, v...)
			_, _ = fmt.Fprintf(buffer, "%s\n", formattedMapTree)
		default:
			fmt.Fprintf(buffer, "%s%s:%-*s %v\n", indent, node.Key, maxKeyLength-len(node.Key)+2, "", node.Value)
		}
	}

	return buffer.String()

}
