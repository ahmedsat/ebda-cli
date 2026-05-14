package utils_test

import (
	"testing"

	"github.com/ahmedsat/ebda-cli/utils"
)

func TestAssertTrueDoesNotExit(t *testing.T) {
	utils.Assert(true, "should not exit")
}
