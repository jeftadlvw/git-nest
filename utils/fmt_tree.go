package utils

import (
	"bytes"
	"fmt"
	"strings"
)

/*
Node is a small key-value structure to fake an ordered dictionary.
Should only be used for formatting with FmtTree.
*/
type Node struct {
	Key   string
	Value interface{}
}

/*
FmtTree formats an array of Node to output an ordered information tree.
Allows nesting of Node.
*/
func FmtTree(indent string, rootLevel bool, nodes ...Node) string {

	buffer := bytes.NewBufferString("")

	maxKeyLength := 0
	for _, node := range nodes {
		maxKeyLength = max(maxKeyLength, len(node.Key))
	}

	localIndent := indent
	if rootLevel {
		localIndent = ""
	}

	for _, node := range nodes {
		switch v := node.Value.(type) {
		case []Node:
			_, _ = fmt.Fprintf(buffer, "%s%s:\n", localIndent, node.Key)
			formattedMapTree := FmtTree(localIndent+indent, false, v...)
			_, _ = fmt.Fprintf(buffer, "%s\n", formattedMapTree)
		default:
			_, _ = fmt.Fprintf(buffer, "%s%s:%-*s%v\n", localIndent, node.Key, maxKeyLength-len(node.Key)+3, "", node.Value)
		}
	}

	return strings.TrimSuffix(buffer.String(), "\n")
}
