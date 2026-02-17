//go:build !release

package cgo

/*
#include "lualib.h"
#include "lauxlib.h"
#include "custom_cgo_helper.h"
*/
import "C"

func (s *State) libs() {
	// s.PushCFunction(CFunction(C.GoFollowUp))
	// s.SetGlobal("FollowUp")
}
