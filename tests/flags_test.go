package tests

import (
	"os"
	"testing"

	"github.com/chapgx/rhombifer"
)

func TestFlags(t *testing.T) {
	r := rhombifer.Root()
	cmd := &rhombifer.Command{Name: "cmd", ShortDesc: "im a command", Run: func(args ...string) error { return nil }}
	cmd.AddFlags(
		rhombifer.NewFlag("foo", "im foo").SetRequired().SetShortFormat("f"),
		rhombifer.NewFlag("bar", "im bar").SetShortFormat("b").SetValuesRequired(),
	)
	r.AddSubs(cmd)

	os.Args = []string{"programanme", "cmd", "--foo", "hello", "world"}

	e := rhombifer.Start()
	if e != nil {
		t.Fatal(e)
	}

	expectedvalues := []string{"hello", "world"}

	// values should defer to command if flags does not require it
	for idx, expectedval := range expectedvalues {
		got := cmd.Values[idx]
		if got != expectedval {
			t.Fatalf("expected %q got %q", expectedval, got)
		}
	}

	os.Args = []string{"programanme", "cmd", "--bar", "hello", "world"}

	e = rhombifer.Start()
	if e != nil {
		t.Fatal(e)
	}

	barflag := cmd.CheckForFlag("bar")
	// values should be kept in the flag if it requires it
	for idx, expectedval := range expectedvalues {
		got := barflag.Values[idx]
		if got != expectedval {
			t.Fatalf("expected %q got %q", expectedval, got)
		}
	}

	// NOTE: testing with short flag
	os.Args = []string{"programanme", "cmd", "-b", "hello", "world"}

	e = rhombifer.Start()
	if e != nil {
		t.Fatal(e)
	}

	// values should be kept in the flag if it requires it
	for idx, expectedval := range expectedvalues {
		// BUG: breaks here  this is probably a parser error where short flags do not get values assign
		got := barflag.Values[idx]
		if got != expectedval {
			t.Fatalf("expected %q got %q", expectedval, got)
		}
	}
}
