package tests

import (
	"fmt"
	"github.com/jeftadlvw/git-nest/utils"
	"testing"
)

func TestFmtTree(t *testing.T) {

	cases := []struct {
		config   utils.FmtTreeConfig
		tree     []utils.Node
		expected string
	}{
		{utils.FmtTreeConfig{Indent: ""}, []utils.Node{}, ""},
		{utils.FmtTreeConfig{Indent: "  "}, []utils.Node{}, ""},
		{utils.FmtTreeConfig{Indent: ""}, []utils.Node{{"Foo", "foo"}, {"Bar", 123}}, "Foo:   foo\nBar:   123"},
		{utils.FmtTreeConfig{Indent: ""}, []utils.Node{{"Foo", "foo"}, {"BarLong", 123}}, "Foo:       foo\nBarLong:   123"},
		{utils.FmtTreeConfig{Indent: "  "}, []utils.Node{{"Foo", "foo"}, {"Bar", 123}}, "Foo:   foo\nBar:   123"},
		{utils.FmtTreeConfig{Indent: "  ", NewLinesAtTopLevel: true}, []utils.Node{{"Foo", "foo"}, {"Bar", 123}}, "Foo:   foo\n\nBar:   123"},
		{utils.FmtTreeConfig{Indent: ""}, []utils.Node{{"Foo", []int{1, 2, 3}}}, "Foo:   [1 2 3]"},
		{utils.FmtTreeConfig{Indent: ""}, []utils.Node{{"Foo", []interface{}{1, 2, "hi"}}}, "Foo:   [1 2 hi]"},
		{utils.FmtTreeConfig{Indent: "   "}, []utils.Node{{"Foo", []utils.Node{{"Bar", 0.34}, {"Baz", "baz"}}}, {"Noc", "noc"}}, "Foo:\n   Bar:   0.34\n   Baz:   baz\nNoc:   noc"},
		{utils.FmtTreeConfig{Indent: "   ", NewLinesAtTopLevel: true}, []utils.Node{{"Foo", []utils.Node{{"Bar", 0.34}, {"Baz", "baz"}}}, {"Noc", "noc"}}, "Foo:\n   Bar:   0.34\n   Baz:   baz\n\nNoc:   noc"},
	}

	for index, tc := range cases {
		t.Run(fmt.Sprintf("TestFmtTree-%d", index+1), func(t *testing.T) {
			output := utils.FmtTree(tc.config, tc.tree...)
			if output != tc.expected {
				t.Fatalf("unexpected results:\nExpected:\n>%s<\n\nGot:\n>%s<", tc.expected, output)
			}
		})
	}
}
