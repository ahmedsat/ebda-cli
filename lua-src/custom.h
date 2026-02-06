#ifndef custom_cgo_helper_h
#define custom_cgo_helper_h

#include "lualib.h"
#include "lauxlib.h"

int GoCall(lua_State *L);

int go_caller(lua_State *L) {return GoCall(L);}

#endif /* custom_cgo_helper_h */