package templater

import (
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func hasUnimplementedPureVirtualFunctions(className string) bool {
	for _, f := range parser.State.ClassMap[className].Functions {

		if !f.Checked {
			cppFunction(f)
			goFunction(f)
			f.Checked = true
		}

		if f.Virtual == parser.PURE && !f.IsSupported() {
			return true
		}
	}
	return false
}

func goModule(module string) string {
	return strings.ToLower(strings.TrimPrefix(module, "Qt"))
}

func UseStub(force bool, module string, mode int) bool {
	return force || (utils.QT_STUB() && mode == NONE && !(module == "QtAndroidExtras" || module == "QtSailfish"))
}

func UseJs() bool { return parser.UseJs() } //TODO: remove

func buildTags(module string, stub bool, mode int, tags string) string {
	switch {
	case stub:
		{
			if module == "QtAndroidExtras" || module == "androidextras" {
				return "// +build !android,!android_emulator"
			}
			return "// +build !sailfish,!sailfish_emulator"
		}

	case mode == MINIMAL:
		{
			return "// +build minimal"
		}

	case mode == MOC:
		{
			if tags != "" {
				return "// +build " + tags
			}
			return ""
		}

	case module == "QtAndroidExtras", module == "androidextras":
		{
			return "// +build android android_emulator"
		}

	case module == "QtSailfish", module == "sailfish":
		{
			return "// +build sailfish sailfish_emulator"
		}

	default:
		{
			return "// +build !minimal"
		}
	}
}
