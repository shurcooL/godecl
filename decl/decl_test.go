package decl_test

import (
	"testing"

	"github.com/shurcooL/play/228/godecl/decl"
)

func TestGoToEnglish(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			"var x int",
			"declare variable x as int",
		},
		{
			"var x, y int",
			"declare variables x and y as int",
		},
		{
			"var x int = 1",
			"declare variable x as int with initial value 1",
		},
		{
			"var x, y int = 1, 2",
			"declare variables x and y as int with initial values 1 and 2",
		},
		{
			"var x = 1",
			"declare variable x with initial value 1",
		},
		{
			"var x, y = 1, 2",
			"declare variables x and y with initial values 1 and 2",
		},
		{
			"var (x int; y string)",
			"declare variable x as int, variable y as string",
		},
		{
			"var (x, y int; a, b string)",
			"declare variables x and y as int, variables a and b as string",
		},
	}
	for _, tc := range tests {
		got, err := decl.GoToEnglish(tc.in)
		if err != nil {
			t.Error("got error:", err)
			continue
		}
		if got != tc.want {
			t.Errorf("got: %q, want: %q", got, tc.want)
		}
	}
}
