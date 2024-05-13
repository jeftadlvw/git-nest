package utils

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
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
FmtTreeConfig is a configuration structure when formatting a Node tree.
*/
type FmtTreeConfig struct {
	/*
		Indent defines the indent contents for subnodes.
	*/
	Indent string

	/*
		NewLinesAtTopLevel defines if there should be any spacing between top-level nodes.
	*/
	NewLinesAtTopLevel bool
}

/*
FmtTree formats an array of Node to output an ordered information tree.
Allows nesting of Node.
*/
func FmtTree(config FmtTreeConfig, nodes ...Node) string {
	return fmtTree(config.Indent, true, config.NewLinesAtTopLevel, nodes...)
}

func fmtTree(indent string, rootNodes bool, newLinesAtTopLevel bool, nodes ...Node) string {

	buffer := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(buffer, 0, 0, 3, ' ', tabwriter.TabIndent)

	maxKeyLength := 0
	for _, node := range nodes {
		maxKeyLength = max(maxKeyLength, len(node.Key))
	}

	localIndent := indent
	if rootNodes {
		localIndent = ""
	}

	for index, node := range nodes {
		switch v := node.Value.(type) {
		case []Node:
			_, _ = fmt.Fprintf(tabWriter, "%s%s:\n", localIndent, node.Key)
			formattedMapTree := fmtTree(localIndent+indent, false, false, v...)
			_, _ = fmt.Fprintf(tabWriter, "%s\n", formattedMapTree)
		default:
			_, _ = fmt.Fprintf(tabWriter, "%s%s:\t%v\n", localIndent, node.Key, node.Value)
		}

		if newLinesAtTopLevel && index != len(nodes)-1 {
			_, _ = fmt.Fprintln(tabWriter, "")
		}
	}

	_ = tabWriter.Flush()
	return strings.TrimSuffix(buffer.String(), "\n")
}
