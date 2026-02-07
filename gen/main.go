package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	sbc := strings.Builder{}
	sbgo := strings.Builder{}

	sbc.WriteString("#ifndef custom_cgo_helper_h\n")
	sbc.WriteString("#define custom_cgo_helper_h\n")
	sbc.WriteString("#include \"lualib.h\"\n")
	sbc.WriteString("#include \"lauxlib.h\"\n")

	sbgo.WriteString("package cgo\n\n")

	sbgo.WriteString("/*\n")
	sbgo.WriteString("#include \"lualib.h\"\n")
	sbgo.WriteString("#include \"lauxlib.h\"\n")
	sbgo.WriteString("#include \"custom_cgo_helper.h\"\n")
	sbgo.WriteString("*/\n")
	sbgo.WriteString("import \"C\"\n\n")
	sbgo.WriteString("func (s *State) libs() {\n")

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		bytes, err := os.ReadFile(path)
		lines := strings.SplitSeq(string(bytes), "\n")

		for line := range lines {
			if !strings.HasPrefix(line, "//export Go") {
				continue
			}
			funcName := strings.TrimPrefix(line, "//export Go")
			fmt.Fprintf(&sbc, "int Go%s(lua_State *L);\n", funcName)
			fmt.Fprintf(&sbgo, "s.PushCFunction(CFunction(C.Go%s))\n", funcName)
			fmt.Fprintf(&sbgo, "s.SetGlobal(%q)\n", funcName)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	sbc.WriteString("#endif // custom_cgo_helper_h\n")
	sbgo.WriteString("}\n")

	err = os.WriteFile("cgo/custom_cgo_helper.h", []byte(sbc.String()), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("cgo/custom_cgo_helper.go", []byte(sbgo.String()), 0644)
	if err != nil {
		panic(err)
	}

}
