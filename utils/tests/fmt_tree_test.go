package tests

import (
	"github.com/jeftadlvw/git-nest/utils"
	"testing"
)

func TestFmtTree(t *testing.T) {

	cases := []struct {
		indent   string
		tree     []utils.Node
		expected string
	}{
		{"", []utils.Node{}, ""},
		{"  ", []utils.Node{}, ""},
		{"", []utils.Node{{"Foo", "foo"}, {"Bar", 123}}, "Foo:   foo\nBar:   123"},
		{"", []utils.Node{{"Foo", "foo"}, {"BarLong", 123}}, "Foo:       foo\nBarLong:   123"},
		{"  ", []utils.Node{{"Foo", "foo"}, {"Bar", 123}}, "Foo:   foo\nBar:   123"},
		{"", []utils.Node{{"Foo", []int{1, 2, 3}}}, "Foo:   [1 2 3]"},
		{"", []utils.Node{{"Foo", []interface{}{1, 2, "hi"}}}, "Foo:   [1 2 hi]"},
		{"   ", []utils.Node{{"Foo", []utils.Node{{"Bar", 0.34}, {"Baz", "baz"}}}}, "Foo:\n   Bar:   0.34\n   Baz:   baz"},
	}

	for _, tc := range cases {
		output := utils.FmtTree(tc.indent, true, tc.tree...)
		if output != tc.expected {
			t.Errorf("FmtTree() for %v returned unexpected results:\nExpected:\n>%s<\n\nGot:\n>%s<", tc.indent, tc.expected, output)
		}
	}

}
