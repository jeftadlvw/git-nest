package tests

import (
	"github.com/jeftadlvw/git-nest/utils"
	"testing"
)

func TestStringInsert(t *testing.T) {

	// empty start and end deliminator
	_, err := utils.StringInsert("", "", "", "")
	if err == nil {
		t.Fatalf("StringInsert() with empty start should return error, but did not")
	}

	_, err = utils.StringInsert("", "", "foo", "")
	if err == nil {
		t.Fatalf("StringInsert() with empty end should return error, but did not")
	}

	const startDelimiter = "@@start"
	const endDelimiter = "@@end"
	const insertString = "Hakuna Matata"
	const originalString = `Lorem ipsum dolor sit amet

@start

@@start
replace here

@@end

@end

Lorem ipsum dolor sit amet
`
	const targetString = `Lorem ipsum dolor sit amet

@start

Hakuna Matata

@end

Lorem ipsum dolor sit amet
`

	// non-existing start and end deliminator
	_, err = utils.StringInsert(originalString, "", "##start", endDelimiter)
	if err == nil {
		t.Fatalf("StringInsert() with non-existings start should return error, but did not")
	}

	_, err = utils.StringInsert(originalString, "", startDelimiter, "###end")
	if err == nil {
		t.Fatalf("StringInsert() with non-existings end should return error, but did not")
	}

	// start and end deliminator that occur often
	_, err = utils.StringInsert(originalString, "", "@start", endDelimiter)
	if err == nil {
		t.Fatalf("StringInsert() with recurring start should return error, but did not")
	}

	_, err = utils.StringInsert(originalString, "", startDelimiter, "@start")
	if err == nil {
		t.Fatalf("StringInsert() with recurring end should return error, but did not")
	}

	output, err := utils.StringInsert(originalString, insertString, startDelimiter, endDelimiter)
	if err != nil {
		t.Fatalf("StringInsert() should not return error, but did: %s", err)
	}

	if output != targetString {
		t.Fatalf("StringInsert() returned unexpected results:\nExpected:\n>%s<\n\nGot:\n>%s<", targetString, output)
	}
}

func TestStringInsertAtFirst(t *testing.T) {

	// empty start and end deliminator
	_, err := utils.StringInsertAtFirst("", "", "", "")
	if err == nil {
		t.Fatalf("StringInsert() with empty start should return error, but did not")
	}

	_, err = utils.StringInsertAtFirst("", "", "foo", "")
	if err == nil {
		t.Fatalf("StringInsert() with empty end should return error, but did not")
	}

	const startDelimiter = "@@start"
	const endDelimiter = "@@end"
	const insertString = "Hakuna Matata"
	const originalString = `Lorem ipsum dolor sit amet

@@start

@@start
@@start
replace here

@@end

@@end

@@start
Some ignored stuff
@@end

Lorem ipsum dolor sit amet
`
	const targetString = `Lorem ipsum dolor sit amet

Hakuna Matata

@@end

@@start
Some ignored stuff
@@end

Lorem ipsum dolor sit amet
`

	// non-existing start and end deliminator
	_, err = utils.StringInsertAtFirst(originalString, "", "##start", endDelimiter)
	if err == nil {
		t.Fatalf("StringInsert() with non-existings start should return error, but did not")
	}

	_, err = utils.StringInsertAtFirst(originalString, "", startDelimiter, "###end")
	if err == nil {
		t.Fatalf("StringInsert() with non-existings end should return error, but did not")
	}

	// start and end deliminator that occur often
	output, err := utils.StringInsertAtFirst(originalString, insertString, startDelimiter, endDelimiter)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if output != targetString {
		t.Fatalf("unexpected results:\nExpected:\n>%s<\n\nGot:\n>%s<", targetString, output)
	}
}
