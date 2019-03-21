package converter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
)

func EnumNeedsCppGlue(value string) bool {
	return strings.ContainsAny(value, "()<>~+") || value == "" || value == "0x1FFFFFFFU"
}

func GoEnum(n string, v string, e *parser.Enum) string {
	var _, err = strconv.Atoi(v)
	switch {
	case EnumNeedsCppGlue(v):
		{
			e.NoConst = true
			if parser.UseJs() {
				if parser.UseWasm() {
					return fmt.Sprintf("int64(qt.WASM.Call(\"_%v_%v_Type\").Int())", strings.Split(e.Fullname, "::")[0], n)
				}
				return fmt.Sprintf("qt.WASM.Call(\"_%v_%v_Type\").Int64()", strings.Split(e.Fullname, "::")[0], n)
			}
			return fmt.Sprintf("C.%v_%v_Type()", strings.Split(e.Fullname, "::")[0], n)
		}

	case strings.Contains(v, "0x"):
		{
			return v
		}

	case err != nil:
		{
			if c, ok := parser.State.ClassMap[class(goEnum(e, v))]; ok && module(c.Module) != module(e) && module(c.Module) != "" {
				return fmt.Sprintf("%v.%v", module(c.Module), goEnum(e, v))
			}
			return goEnum(e, v)
		}
	}

	return v
}
