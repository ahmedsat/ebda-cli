package utils_test

import (
	"strings"
	"testing"

	"github.com/ahmedsat/ebda-cli/utils"
)

func TestKmlColor(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "red", in: "#FF0000", want: "7d0000ff"},
		{name: "without hash", in: "AABBCC", want: "7dccbbaa"},
		{name: "black", in: "#000000", want: "7d000000"},
		{name: "white", in: "#FFFFFF", want: "7dffffff"},
		{name: "invalid", in: "bad", want: "7d00ff00"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := strings.ToLower(utils.KmlColor(tc.in)); got != tc.want {
				t.Fatalf("KmlColor(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
