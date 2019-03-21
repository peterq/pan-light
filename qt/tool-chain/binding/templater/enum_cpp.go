package templater

import (
	"fmt"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
)

func cppEnum(e *parser.Enum, v *parser.Value) string {
	o := fmt.Sprintf("%v\n{\n\t%v\n}", cppEnumHeader(e, v), cppEnumBody(e, v))
	if UseJs() {
		o = "EMSCRIPTEN_KEEPALIVE\n" + o
	}
	return o
}

func cppEnumHeader(enum *parser.Enum, value *parser.Value) string {
	return fmt.Sprintf("int %v_%v_Type()", enum.ClassName(), value.Name)
}

func cppEnumBody(enum *parser.Enum, value *parser.Value) string {
	//TODO: check for "since" tag in enums

	//needed for sailfish with 5.6 docs
	if strings.HasPrefix(value.Name, "MV_") || strings.HasPrefix(value.Name, "PM_") ||
		strings.HasPrefix(value.Name, "SH_") || strings.HasPrefix(value.Name, "ISODate") ||
		strings.HasPrefix(value.Name, "TlsV1_") {
		return fmt.Sprintf(`#if QT_VERSION >= 0x056000
		return %v::%v;
	#else
		return 0;
	#endif`, enum.ClassName(), value.Name)
	}

	//needed for msys2 with 5.7 docs
	if strings.HasPrefix(value.Name, "PE_") || strings.HasPrefix(value.Name, "SE_") {
		return fmt.Sprintf(`#if QT_VERSION >= 0x057000
		return %v::%v;
	#else
		return 0;
	#endif`, enum.ClassName(), value.Name)
	}

	return fmt.Sprintf("return %v::%v;", enum.ClassName(), value.Name)
}
