//go:build !release

package cgo

/*
#cgo LDFLAGS: -L../lua-src -llua -lm
#cgo CFLAGS: -I../lua-src
#include <stdlib.h>
#include "lualib.h"
#include "lauxlib.h"
#include "custom_cgo_helper.h"
*/
import "C"
import (
	"fmt"
)

const MULTRET = -1

type KContext C.lua_KContext
type KFunction C.lua_KFunction
type CFunction C.lua_CFunction
type State C.lua_State

func (s *State) Close() {
	C.lua_close((*C.lua_State)(s))
}

func (s *State) DoFile(filename string) error {
	err := s.LoadFile(filename)
	if err != nil {
		return err
	}

	err = s.Pcall(0, MULTRET, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s *State) Pcall(na, nr, errF int) error {
	return s.PcallK(na, nr, errF, 0, nil)
}

func (s *State) PcallK(na, nr, errF int, ctx KContext, k KFunction) error {
	if C.lua_pcallk((*C.lua_State)(s), C.int(na), C.int(nr), C.int(errF), C.lua_KContext(ctx), C.lua_KFunction(k)) != 0 {
		return s.GetError()
	}
	return nil
}

func (s *State) Loadfilex(filename string, mode string) error {

	CFilename, free := ToCString(filename)
	CMode, free2 := ToCString(mode)
	defer free()
	defer free2()

	if C.luaL_loadfilex((*C.lua_State)(s), CFilename, CMode) != 0 {
		return s.GetError()
	} else {
		return nil
	}
}

func (s *State) LoadFile(filename string) error {
	return s.Loadfilex(filename, "")
}

func NewLuaState() (*State, error) {
	s := C.luaL_newstate()
	if s == nil {
		return nil, fmt.Errorf("failed to create lua state")
	}
	return (*State)(s), nil

}

func (s *State) GetError() error {
	return fmt.Errorf("%s", C.GoString(C.luaL_tolstring((*C.lua_State)(s), -1, nil)))
}

// tolstring (lua_State *L, int idx, size_t *len)
func (s *State) ToString(idx int) string {
	return C.GoString(C.luaL_tolstring((*C.lua_State)(s), C.int(idx), nil))
}

func (s *State) OpenLibs() {
	s.OpenSelectedLibs(^0, 0)
	s.libs()
}

func (s *State) SetGlobal(name string) {
	cName, free := ToCString(name)
	defer free()
	C.lua_setglobal((*C.lua_State)(s), cName)
}

func (s *State) PushCFunction(fn CFunction) {
	s.PushCClosure(fn, 0)
}

func (s *State) PushCClosure(fn CFunction, n int) {
	C.lua_pushcclosure((*C.lua_State)(s), fn, C.int(n))
}

func (s *State) OpenSelectedLibs(load, preload int) {
	C.luaL_openselectedlibs((*C.lua_State)(s), C.int(load), C.int(preload))
}
