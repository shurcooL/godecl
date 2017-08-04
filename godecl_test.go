package main_test

import "testing"

// TODO: Move to library package.
func Test(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"var x int", "declare variables x as int"},
		{"var x, y int", "declare variables x and y as int"},
		{"", ""},
		{"", ""},
		{"", ""},
		{"var (x int; y, z int)", "declare variable x as int and y and z as ints"},
		{"", ""},
		{"var (a, b int; x, y string)", ""},
		{"", ""},
	}
	for _, tc := range tests {
		got, err := GoToEnglish(tc.in)
		if err != nil {
			t.Error("got error:", err)
			continue
		}
		if got != tc.want {
			t.Errorf("got: %q, want: %q", got, tc.want)
		}
	}
}
