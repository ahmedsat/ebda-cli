package utils_test

import (
	"strings"
	"testing"
	"time"

	"github.com/ahmedsat/ebda-cli/utils"
)

func TestTimeLayout(t *testing.T) {
	parsed, err := time.Parse(utils.TimeLayout, "1-6-2023")
	if err != nil {
		t.Fatalf("parse valid date: %v", err)
	}
	if got := parsed.Format(utils.TimeLayout); got != "1-6-2023" {
		t.Fatalf("round trip = %q, want %q", got, "1-6-2023")
	}
	if _, err := time.Parse(utils.TimeLayout, "2023-06-01"); err == nil {
		t.Fatal("expected ISO date to fail")
	}
}

func TestFilter(t *testing.T) {
	in := []int{5, 3, 8, 1, 9}
	got := utils.Filter(in, func(n int) bool { return n > 4 })
	want := []int{5, 8, 9}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}

	got = utils.Filter[int](nil, func(n int) bool { return true })
	if len(got) != 0 {
		t.Fatalf("nil slice result len = %d, want 0", len(got))
	}
}

func TestSameAfterSanitize(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{name: "identical", a: "hello", b: "hello", want: true},
		{name: "case", a: "Hello", b: "hello", want: true},
		{name: "space", a: "  hello  ", b: "he llo", want: true},
		{name: "different", a: "apple", b: "orange", want: false},
		{name: "empty", a: "", b: "", want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := utils.SameAfterSanitize(tc.a, tc.b); got != tc.want {
				t.Fatalf("SameAfterSanitize(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{in: "hello world", want: "HelloWorld"},
		{in: "hello", want: "Hello"},
		{in: "HelloWorld", want: "HelloWorld"},
		{in: "hello_world", want: "HelloWorld"},
		{in: "HELLO WORLD", want: "HelloWorld"},
		{in: "", want: ""},
	}

	for _, tc := range cases {
		t.Run(strings.ReplaceAll(tc.in, " ", "_"), func(t *testing.T) {
			if got := utils.ToPascalCase(tc.in); got != tc.want {
				t.Fatalf("ToPascalCase(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
