package templater

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func HTemplate(m string, mode int, tags string) []byte {
	utils.Log.WithField("module", m).Debug("generating h")

	var bb = new(bytes.Buffer)
	defer bb.Reset()

	if m != parser.MOC {
		m = "Qt" + m
	}

	//header
	fmt.Fprintf(bb, "%v\n\n", buildTags(m, false, mode, tags))

	fmt.Fprint(bb, "#pragma once\n\n")

	var hash string
	if m == parser.MOC {
		hash = "_" + parser.SortedClassesForModule(m, true)[0].Hash() //TODO:
	}
	fmt.Fprintf(bb, "#ifndef GO_%v%v_H\n", strings.ToUpper(m), hash)
	fmt.Fprintf(bb, "#define GO_%v%v_H\n\n", strings.ToUpper(m), hash)

	fmt.Fprint(bb, "#include <stdint.h>\n\n")

	fmt.Fprint(bb, "#ifdef __cplusplus\n")
	for _, c := range parser.SortedClassNamesForModule(m, true) {
		if parser.State.ClassMap[c].IsSubClassOfQObject() {
			if m == parser.MOC {
				fmt.Fprintf(bb, "class %v;\n", c)
				fmt.Fprintf(bb, "void %[1]v_%[1]v_QRegisterMetaTypes();\n", c)
			} else {
				fmt.Fprintf(bb, "int %[1]v_%[1]v_QRegisterMetaType();\n", c)
			}
		}
	}

	fmt.Fprint(bb, "extern \"C\" {\n#endif\n\n")

	if !UseJs() {
		fmt.Fprintf(bb, "struct %v_PackedString { char* data; long long len; };\n", strings.Title(m))
		fmt.Fprintf(bb, "struct %v_PackedList { void* data; long long len; };\n", strings.Title(m))
	}

	//body
	for _, c := range parser.SortedClassesForModule(m, true) {
		cTemplate(bb, c, cppEnumHeader, cppFunctionHeader, ";\n", false)
	}

	//footer
	fmt.Fprint(bb, "\n#ifdef __cplusplus\n}\n#endif\n\n#endif")

	//TODO: regexp
	if mode == MOC {
		pre := bb.String()
		bb.Reset()
		libsm := make(map[string]struct{}, 0)
		for _, c := range parser.State.ClassMap {
			if c.Pkg != "" && c.IsSubClassOfQObject() {
				libsm[c.Module] = struct{}{}
			}
		}

		var libs []string
		for k := range libsm {
			libs = append(libs, k)
		}
		libs = append(libs, m)

		for _, c := range parser.SortedClassesForModule(strings.Join(libs, ","), true) {
			hName := c.Hash()
			sep := []string{"\"_", "LIVE_", " ", "\t", "\n", "\r", "(", ")", ":", ";", "*", "<", ">", "&", "~", "{", "}", "[", "]", "_", "callback"}
			for _, p := range sep {
				for _, s := range sep {
					if s == "callback" {
						continue
					}
					pre = strings.Replace(pre, p+c.Name+s, p+c.Name+hName+s, -1)
				}
			}
		}
		bb.WriteString(pre)
	}

	if !UseJs() {
		return bb.Bytes()
	}
	tmp := bb.String()
	for _, l := range strings.Split(tmp, "\n") {
		if strings.Contains(l, "emscripten::val") {
			tmp = strings.Replace(tmp, l, "", -1)
		}
	}
	return []byte(tmp)
}
