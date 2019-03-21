package converter

//TODO: GLchar, GLbyte

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
)

func GoOutput(name, value string, f *parser.Function, p string) string {
	return goOutput(name, value, f, p)
}
func goOutput(name, value string, f *parser.Function, p string) string {
	vOld := value

	name = parser.CleanName(name, value)
	value = parser.CleanValue(value)

	switch value {
	case "char", "qint8", "uchar", "quint8", "GLubyte", "QString":
		{
			if !parser.UseJs() { //TODO: support []byte in js as well
				switch value {
				case "char", "qint8", "uchar", "quint8", "GLubyte":
					if len(f.Parameters) <= 4 &&
						(strings.Contains(strings.ToLower(f.Name), "read") ||
							strings.Contains(strings.ToLower(f.Name), "write") ||
							strings.Contains(strings.ToLower(f.Name), "data")) {
						for _, p := range f.Parameters {
							if strings.Contains(p.Value, "int") && f.Parameters[0].Value == vOld {
								return fmt.Sprintf("cGoUnpackBytes(%v)", name)
							}
						}
					}
				}
			}

			return func() string {
				var out = fmt.Sprintf("cGoUnpackString(%v)", name)
				if strings.Contains(p, "error") {
					return fmt.Sprintf("errors.New(%v)", out)
				}
				return out
			}()
		}

	case "QStringList":
		{
			return fmt.Sprintf("strings.Split(cGoUnpackString(%v), \"|\")", name)
		}

	case "void", "GLvoid", "":
		{
			return name
		}

	case "bool", "GLboolean":
		{
			return fmt.Sprintf("int8(%v) != 0", name)
		}

	case "short", "qint16", "GLshort":
		{
			return fmt.Sprintf("int16(%v)", name)
		}

	case "ushort", "unsigned short", "quint16", "GLushort":
		{
			return fmt.Sprintf("uint16(%v)", name)
		}

	case "int", "qint32", "GLint", "GLsizei", "GLintptrARB", "GLsizeiptrARB", "GLfixed", "GLclampx":
		{
			return fmt.Sprintf("int(int32(%v))", name)
		}

	case "uint", "unsigned int", "quint32", "GLenum", "GLbitfield", "GLuint", "QRgb":
		{
			return fmt.Sprintf("uint(uint32(%v))", name)
		}

	case "long":
		{
			return fmt.Sprintf("int(int32(%v))", name)
		}

	case "ulong", "unsigned long":
		{
			return fmt.Sprintf("uint(uint32(%v))", name)
		}

	case "longlong", "long long", "qlonglong", "qint64":
		{
			return fmt.Sprintf("int64(%v)", name)
		}

	case "ulonglong", "unsigned long long", "qulonglong", "quint64":
		{
			return fmt.Sprintf("uint64(%v)", name)
		}

	case "float", "GLfloat", "GLclampf":
		{
			return fmt.Sprintf("float32(%v)", name)
		}

	case "double", "qreal":
		{
			return fmt.Sprintf("float64(%v)", name)
		}

	case "uintptr_t", "uintptr", "quintptr", "WId":
		{
			return fmt.Sprintf("uintptr(%v)", name)
		}

		//non std types

	case "T", "JavaVM", "jclass", "jobject":
		{
			switch f.TemplateModeJNI {
			case "Boolean":
				{
					return fmt.Sprintf("int8(%v) != 0", name)
				}

			case "Int":
				{
					return fmt.Sprintf("int(int32(%v))", name)
				}

			case "Void":
				{
					return name
				}
			}

			return fmt.Sprintf("unsafe.Pointer(%v)", name)
		}
	}

	switch {
	case isEnum(f.ClassName(), value):
		{
			if c, ok := parser.State.ClassMap[class(cppEnum(f, value, false))]; ok && module(c.Module) != module(f) && module(c.Module) != "" {
				if _, ok := parser.State.ClassMap[f.ClassName()].WeakLink[c.Module]; ok {
					return fmt.Sprintf("int64(%v)", name)
				}
				return fmt.Sprintf("%v.%v(%v)", module(c.Module), goEnum(f, value), name)
			}
			return fmt.Sprintf("%v(%v)", goEnum(f, value), name)
		}

	case isClass(value):
		{
			if strings.Contains(value, ".") {
				value = strings.Split(value, ".")[1]
			}
			if m := module(parser.State.ClassMap[value].Module); m != module(f) {
				if _, ok := parser.State.ClassMap[f.ClassName()].WeakLink[parser.State.ClassMap[value].Module]; ok {
					return fmt.Sprintf("unsafe.Pointer(%v)", name)
				}
				return fmt.Sprintf("%v.New%vFromPointer(%v)", m, strings.Title(value), name)
			}
			return fmt.Sprintf("New%vFromPointer(%v)", strings.Title(value), name)
		}

	case parser.IsPackedList(value):
		{
			return fmt.Sprintf("func(l C.struct_%v_PackedList)%v{out := make(%v, int(l.len))\ntmpList := New%vFromPointer(l.data)\nfor i:=0;i<len(out);i++{ out[i] = tmpList.__%v_atList%v(i) }\nreturn out}(%v)", strings.Title(parser.State.ClassMap[f.ClassName()].Module), goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, f.OverloadNumber, name)
		}

	case parser.IsPackedMap(value):
		{
			return fmt.Sprintf("func(l C.struct_%v_PackedList)%v{out := make(%v, int(l.len))\ntmpList := New%vFromPointer(l.data)\nfor i,v:=range tmpList.__%v_keyList(){ out[v] = tmpList.__%v_atList%v(v, i) }\nreturn out}(%v)", strings.Title(parser.State.ClassMap[f.ClassName()].Module), goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, f.Name, f.OverloadNumber, name)
		}
	}

	f.Access = fmt.Sprintf("unsupported_goOutput(%v)", value)
	return f.Access
}

func GoOutputFailed(value string, f *parser.Function, p string) string {
	return goOutputFailed(value, f, p)
}
func goOutputFailed(value string, f *parser.Function, p string) string {
	var vOld = value

	value = parser.CleanValue(value)

	switch value {
	case "char", "qint8", "uchar", "quint8", "GLubyte", "QString":
		{
			if strings.Contains(p, "error") {
				return "errors.New(\"\")"
			}
			return "\"\""
		}

	case "QStringList":
		{
			return "make([]string, 0)"
		}

	case "void", "GLvoid", "":
		{
			if strings.Contains(vOld, "*") {
				return "nil"
			}

			return ""
		}

	case "bool", "GLboolean":
		{
			return "false"
		}

	case
		"short", "qint16", "GLshort",
		"ushort", "unsigned short", "quint16", "GLushort",

		"int", "qint32", "GLint", "GLsizei", "GLintptrARB", "GLsizeiptrARB", "GLfixed", "GLclampx",
		"uint", "unsigned int", "quint32", "GLenum", "GLbitfield", "GLuint", "QRgb",

		"long",
		"ulong", "unsigned long",

		"longlong", "long long", "qlonglong", "qint64",
		"ulonglong", "unsigned long long", "qulonglong", "quint64",

		"float", "GLfloat", "GLclampf",
		"double", "qreal",

		"uintptr_t", "uintptr", "quintptr", "WId":
		{
			return "0"
		}

		//non std types

	case "T", "JavaVM", "jclass", "jobject":
		{
			switch f.TemplateModeJNI {
			case "Boolean":
				{
					return "false"
				}

			case "Int":
				{
					return "0"
				}

			case "Void":
				{
					return ""
				}
			}

			return "nil"
		}
	}

	switch {
	case isEnum(f.ClassName(), value):
		{
			return "0"
		}

	case isClass(value):
		{
			if strings.Contains(value, ".") {
				value = strings.Split(value, ".")[1]
			}
			if f.TemplateModeJNI == "String" {
				return "\"\""
			}

			return "nil"
		}

	case parser.IsPackedList(value) || parser.IsPackedMap(value):
		{
			return fmt.Sprintf("make(%v, 0)", goType(f, value, p))
		}
	}

	f.Access = fmt.Sprintf("unsupported_goOutputFailed(%v)", value)
	return f.Access
}

func cgoOutput(name, value string, f *parser.Function, p string) string {
	vOld := value

	name = parser.CleanName(name, value)
	value = parser.CleanValue(value)

	switch value {
	case "char", "qint8", "uchar", "quint8", "GLubyte", "QString":
		{
			if !parser.UseJs() { //TODO: support []byte in js as well
				switch value {
				case "char", "qint8", "uchar", "quint8", "GLubyte":
					if len(f.Parameters) <= 4 &&
						(strings.Contains(strings.ToLower(f.Name), "read") ||
							strings.Contains(strings.ToLower(f.Name), "write") ||
							strings.Contains(strings.ToLower(f.Name), "data")) {
						for _, p := range f.Parameters {
							if strings.Contains(p.Value, "int") && f.Parameters[0].Value == vOld {
								return fmt.Sprintf("cGoUnpackBytes(%v)", name)
							}
						}
					}
				}
			}

			out := fmt.Sprintf("cGoUnpackString(%v)", name)
			if parser.UseJs() {
				out = name
			}
			if strings.Contains(p, "error") {
				return fmt.Sprintf("errors.New(%v)", out)
			}
			return out
		}

	case "QStringList":
		{
			if parser.UseJs() {
				return fmt.Sprintf("strings.Split(%v, \"|\")", name)
			}
			return fmt.Sprintf("strings.Split(cGoUnpackString(%v), \"|\")", name)
		}

	case "void", "GLvoid", "":
		{
			return name
		}

	case "bool", "GLboolean":
		{
			if parser.UseJs() {
				return name
			}
			return fmt.Sprintf("int8(%v) != 0", name)
		}

	case "short", "qint16", "GLshort":
		{
			return fmt.Sprintf("int16(%v)", name)
		}

	case "ushort", "unsigned short", "quint16", "GLushort":
		{
			return fmt.Sprintf("uint16(%v)", name)
		}

	case "int", "qint32", "GLint", "GLsizei", "GLintptrARB", "GLsizeiptrARB", "GLfixed", "GLclampx":
		{
			return fmt.Sprintf("int(int32(%v))", name)
		}

	case "uint", "unsigned int", "quint32", "GLenum", "GLbitfield", "GLuint", "QRgb":
		{
			return fmt.Sprintf("uint(uint32(%v))", name)
		}

	case "long":
		{
			return fmt.Sprintf("int(int32(%v))", name)
		}

	case "ulong", "unsigned long":
		{
			return fmt.Sprintf("uint(uint32(%v))", name)
		}

	case "longlong", "long long", "qlonglong", "qint64":
		{
			return fmt.Sprintf("int64(%v)", name)
		}

	case "ulonglong", "unsigned long long", "qulonglong", "quint64":
		{
			return fmt.Sprintf("uint64(%v)", name)
		}

	case "float", "GLfloat", "GLclampf":
		{
			return fmt.Sprintf("float32(%v)", name)
		}

	case "double", "qreal":
		{
			return fmt.Sprintf("float64(%v)", name)
		}

	case "uintptr_t", "uintptr", "quintptr", "WId":
		{
			return fmt.Sprintf("uintptr(%v)", name)
		}

		//non std types

	case "T", "JavaVM", "jclass", "jobject":
		{
			switch f.TemplateModeJNI {
			case "Boolean":
				{
					return fmt.Sprintf("int8(%v) != 0", name)
				}

			case "Int":
				{
					return fmt.Sprintf("int(int32(%v))", name)
				}

			case "Void":
				{
					return name
				}
			}

			return fmt.Sprintf("unsafe.Pointer(%v)", name)
		}
	}

	switch {
	case isEnum(f.ClassName(), value):
		{
			if c, ok := parser.State.ClassMap[class(cppEnum(f, value, false))]; ok && module(c.Module) != module(f) && module(c.Module) != "" {
				if _, ok := parser.State.ClassMap[f.ClassName()].WeakLink[c.Module]; ok {
					return fmt.Sprintf("int64(%v)", name)
				}
				return fmt.Sprintf("%v.%v(%v)", module(c.Module), goEnum(f, value), name)
			}
			return fmt.Sprintf("%v(%v)", goEnum(f, value), name)
		}

	case isClass(value):
		{
			if parser.UseJs() && f.SignalMode != parser.CALLBACK {
				name = fmt.Sprintf("func() uintptr { if %v != js.Undefined { return uintptr(%v.Call(\"Pointer\").Int64()) }; return 0 }()", name, name)
			}

			if strings.Contains(value, ".") {
				value = strings.Split(value, ".")[1]
			}
			if m := module(parser.State.ClassMap[value].Module); m != module(f) {
				if _, ok := parser.State.ClassMap[f.ClassName()].WeakLink[parser.State.ClassMap[value].Module]; ok {
					return fmt.Sprintf("unsafe.Pointer(%v)", name)
				}
				if parser.UseJs() {
					return fmt.Sprintf("%v.New%vFromPointer(unsafe.Pointer(%v))", m, strings.Title(value), name)
				}
				return fmt.Sprintf("%v.New%vFromPointer(%v)", m, strings.Title(value), name)
			}
			if parser.UseJs() {
				return fmt.Sprintf("New%vFromPointer(unsafe.Pointer(%v))", strings.Title(value), name)
			}
			return fmt.Sprintf("New%vFromPointer(%v)", strings.Title(value), name)
		}

	case parser.IsPackedList(value):
		{
			if parser.UseJs() {
				if parser.UseWasm() {
					return fmt.Sprintf("func(l js.Value)%v{out := make(%v, int(l.Get(\"len\").Int()))\ntmpList := New%vFromPointer(unsafe.Pointer(uintptr(l.Get(\"data\").Int())))\nfor i:=0;i<len(out);i++{ out[i] = tmpList.__%v_%v_atList%v(i) }\nreturn out}(%v)", goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, name)
				}
				return fmt.Sprintf("func(l *js.Object)%v{out := make(%v, int(l.Get(\"len\").Int()))\ntmpList := New%vFromPointer(unsafe.Pointer(l.Get(\"data\").Unsafe()))\nfor i:=0;i<len(out);i++{ out[i] = tmpList.__%v_%v_atList%v(i) }\nreturn out}(%v)", goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, name)
			}
			return fmt.Sprintf("func(l C.struct_%v_PackedList)%v{out := make(%v, int(l.len))\ntmpList := New%vFromPointer(l.data)\nfor i:=0;i<len(out);i++{ out[i] = tmpList.__%v_%v_atList%v(i) }\nreturn out}(%v)", strings.Title(parser.State.ClassMap[f.ClassName()].Module), goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, name)
		}

	case parser.IsPackedMap(value):
		{
			if parser.UseJs() {
				if parser.UseWasm() {
					return fmt.Sprintf("func(l js.Value)%v{out := make(%v, int(l.Get(\"len\").Int()))\ntmpList := New%vFromPointer(unsafe.Pointer(uintptr(l.Get(\"data\").Int())))\nfor i,v:=range tmpList.__%v_%v_keyList%v(){ out[v] = tmpList.__%v_%v_atList%v(v, i) }\nreturn out}(%v)", goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, f.Name, name, f.OverloadNumber, name)
				}
				return fmt.Sprintf("func(l *js.Object)%v{out := make(%v, int(l.Get(\"len\").Int()))\ntmpList := New%vFromPointer(unsafe.Pointer(l.Get(\"data\").Unsafe()))\nfor i,v:=range tmpList.__%v_%v_keyList%v(){ out[v] = tmpList.__%v_%v_atList%v(v, i) }\nreturn out}(%v)", goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, f.Name, name, f.OverloadNumber, name)
			}
			return fmt.Sprintf("func(l C.struct_%v_PackedList)%v{out := make(%v, int(l.len))\ntmpList := New%vFromPointer(l.data)\nfor i,v:=range tmpList.__%v_%v_keyList(){ out[v] = tmpList.__%v_%v_atList%v(v, i) }\nreturn out}(%v)", strings.Title(parser.State.ClassMap[f.ClassName()].Module), goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, name, f.Name, name, f.OverloadNumber, name)
		}
	}

	f.Access = fmt.Sprintf("unsupported_cgoOutput(%v)", value)
	return f.Access
}

func CppOutput(name, value string, f *parser.Function) string {
	if strings.HasSuffix(f.Name, "_atList") {
		if parser.UseJs() {
			if f.IsMap {
				out := cppOutput(fmt.Sprintf("({%v tmp = %v->value%v; if (i == %v->size()-1) { %v->~%v(); free(reinterpret_cast<void*>(ptr)); }; tmp; })", value, strings.Split(name, "->")[0], "("+strings.TrimSuffix(strings.Split(name, "_atList(")[1], ", i)")+")", strings.Split(name, "->")[0], strings.Split(name, "->")[0], parser.CleanValue(f.Container)), value, f)
				if !strings.Contains(cppOutput(name, value, f), "emscripten::val") && f.BoundByEmscripten {
					if !strings.Contains(out, "emscripten::val::global") {
						out = "reinterpret_cast<uintptr_t>(" + out + ")"
					}
				}
				return out
			}
			out := cppOutput(fmt.Sprintf("({%v tmp = %v->at%v; if (i == %v->size()-1) { %v->~%v(); free(reinterpret_cast<void*>(ptr)); }; tmp; })", value, strings.Split(name, "->")[0], "("+strings.Split(name, "_atList(")[1], strings.Split(name, "->")[0], strings.Split(name, "->")[0], parser.CleanValue(f.Container)), value, f)
			if !strings.Contains(cppOutput(name, value, f), "emscripten::val") && f.BoundByEmscripten {
				if !strings.Contains(out, "emscripten::val::global") {
					out = "reinterpret_cast<uintptr_t>(" + out + ")"
				}
			}
			return out
		}
		if f.IsMap {
			return cppOutput(fmt.Sprintf("({%v tmp = %v->value%v; if (i == %v->size()-1) { %v->~%v(); free(ptr); }; tmp; })", value, strings.Split(name, "->")[0], "("+strings.TrimSuffix(strings.Split(name, "_atList(")[1], ", i)")+")", strings.Split(name, "->")[0], strings.Split(name, "->")[0], parser.CleanValue(f.Container)), value, f)
		}
		return cppOutput(fmt.Sprintf("({%v tmp = %v->at%v; if (i == %v->size()-1) { %v->~%v(); free(ptr); }; tmp; })", value, strings.Split(name, "->")[0], "("+strings.Split(name, "_atList(")[1], strings.Split(name, "->")[0], strings.Split(name, "->")[0], parser.CleanValue(f.Container)), value, f)
	}
	if strings.HasSuffix(f.Name, "_setList") {
		if len(f.Parameters) == 2 {
			return cppOutput(fmt.Sprintf("%v->insert%v", strings.Split(name, "->")[0], "("+strings.Split(name, "_setList(")[1]), value, f)
		}
		return cppOutput(fmt.Sprintf("%v->append%v", strings.Split(name, "->")[0], "("+strings.Split(name, "_setList(")[1]), value, f)
	}
	if strings.HasSuffix(f.Name, "_newList") {
		return fmt.Sprintf("new %v()", parser.CleanValue(f.Container))
	}
	if strings.HasSuffix(f.Name, "_keyList") {
		return cppOutput(fmt.Sprintf("static_cast<%v*>(ptr)->keys()", f.Container), value, f)
	}
	out := cppOutput(name, value, f)

	if f.BoundByEmscripten && (strings.Contains(CppHeaderOutput(f), "uintptr_t") || strings.Contains(CppHeaderOutput(f), "void*")) && f.SignalMode != parser.CALLBACK {
		return fmt.Sprintf("reinterpret_cast<uintptr_t>(%v)", out)
	}

	if parser.UseJs() && f.SignalMode != parser.CALLBACK {
		for _, p := range f.Parameters {
			if strings.Contains(cppType(f, p.Value), "emscripten::val") && isClass(parser.CleanValue(f.Output)) && !strings.Contains(cppType(f, f.Output), "emscripten::val") && f.BoundByEmscripten {
				if !strings.Contains(out, "emscripten::val::global") && !strings.ContainsAny(out, ";") {
					return fmt.Sprintf("reinterpret_cast<uintptr_t>(%v)", out)
				}
			}
		}
	}

	if parser.UseJs() && f.SignalMode != parser.CALLBACK {
		if isClass(parser.CleanValue(f.Output)) && !strings.Contains(cppType(f, f.Output), "emscripten::val") && f.BoundByEmscripten {
			if !strings.Contains(out, "emscripten::val::global") {
				if strings.Contains(out, "; new") {
					return strings.Replace(strings.Replace(out, "; new", "; reinterpret_cast<uintptr_t>(new", -1), "; })", "); })", -1)
				}
			}
		}
	}

	return out
}

func cppOutputPack(name, value string, f *parser.Function) string {
	var out = CppOutput(name, value, f)

	if strings.Contains(out, "_PackedString") {
		var out = strings.Replace(out, "({ ", "", -1)
		out = strings.Replace(out, " })", "", -1)
		if !strings.HasSuffix(out, ";") {
			out = fmt.Sprintf("%v;", out)
		}
		return strings.Replace(out, "_PackedString", fmt.Sprintf("_PackedString %vPacked =", parser.CleanName(name, value)), -1)
	}

	return ""
}

func cppOutputPacked(name, value string, f *parser.Function) string {
	var out = CppOutput(name, value, f)

	if parser.UseJs() {
		if isClass(parser.CleanValue(value)) && !strings.Contains(out, "emscripten::val::object()") {
			return "reinterpret_cast<uintptr_t>(" + out + ")"
		}
	} else {
		if strings.Contains(out, "_PackedString") {
			return fmt.Sprintf("%vPacked", parser.CleanName(name, value))
		}
	}

	return out
}

//TODO: remove hex encoding once QByteArray <-> ArrayBuffer conversion is possible and/or more TypedArray functions are available for gopherjs/wasm
//TODO: make exemption for QString and QStringList for now? they usually won't need the extra hex encoding ...
//TOOD: or use malloc and simply return a pointer? instead waiting for gopherjs/wasm?
func cppOutputPackingStringForJs(name, length string) string {
	if parser.UseJs() {
		return fmt.Sprintf("emscripten::val ret = emscripten::val::object(); ret.set(\"data\", QByteArray::fromRawData(%v, %v).toHex().toStdString()); ret.set(\"len\", %v); ret;", name, length, length)
	}
	return ""
}

func cppOutputPackingListForJs() string {
	if parser.UseJs() {
		return "emscripten::val ret = emscripten::val::object(); ret.set(\"data\", reinterpret_cast<uintptr_t>(tmpValue)); ret.set(\"len\", tmpValue->size());"
	}
	return ""
}

func cppOutput(name, value string, f *parser.Function) string {
	var vOld = value

	var tHash = sha1.New()
	tHash.Write([]byte(name))
	var tHashName = hex.EncodeToString(tHash.Sum(nil)[:3])

	name = parser.CleanName(name, value)
	value = parser.CleanValue(value)

	switch value {
	case "char", "qint8":
		{
			fSizeVariable := "-1"
			for _, p := range f.Parameters {
				if strings.Contains(p.Value, "int") {
					fSizeVariable = parser.CleanName(p.Name, p.Value)
					break
				}
			}

			if fSizeVariable == "-1" && strings.Contains(strings.ToLower(f.Name), "data") && parser.State.ClassMap[f.ClassName()].HasFunctionWithName("size") {
				fSizeVariable = fmt.Sprintf("static_cast<%v*>(ptr)->size()", f.ClassName())
			}

			if strings.Contains(vOld, "*") {
				if strings.Contains(vOld, "const") {
					if parser.UseJs() {
						return "({ " + cppOutputPackingStringForJs(fmt.Sprintf("const_cast<char*>(%v)", name), fSizeVariable) + " })"
					}
					return fmt.Sprintf("%v_PackedString { const_cast<char*>(%v), %v }", strings.Title(parser.State.ClassMap[f.ClassName()].Module), name, fSizeVariable)
				} else {
					if parser.UseJs() {
						return "({ " + cppOutputPackingStringForJs(name, fSizeVariable) + " })"
					}
					return fmt.Sprintf("%v_PackedString { %v, %v }", strings.Title(parser.State.ClassMap[f.ClassName()].Module), name, fSizeVariable)
				}
			}

			if parser.UseJs() {
				return fmt.Sprintf("({ char t%v = %v; %v })", tHashName, name, cppOutputPackingStringForJs("&t"+tHashName, "-1"))
			}
			return fmt.Sprintf("({ char t%v = %v; %v_PackedString { &t%v, %v }; })", tHashName, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, "-1")
		}

	case "uchar", "quint8", "GLubyte":
		{
			fSizeVariable := "-1"
			if fSizeVariable == "-1" && strings.Contains(strings.ToLower(f.Name), "bits") && len(f.Parameters) == 0 && parser.State.ClassMap[f.ClassName()].HasFunctionWithName("mappedBytes") {
				fSizeVariable = fmt.Sprintf("static_cast<%v*>(ptr)->mappedBytes()", f.ClassName())
			}

			if strings.Contains(vOld, "*") {
				if strings.Contains(vOld, "const") {
					if parser.UseJs() {
						return fmt.Sprintf("({ char* t%v = static_cast<char*>(static_cast<void*>(const_cast<%v*>(%v))); %v })", tHashName, value, name, cppOutputPackingStringForJs("t"+tHashName, fSizeVariable))
					}
					return fmt.Sprintf("({ char* t%v = static_cast<char*>(static_cast<void*>(const_cast<%v*>(%v))); %v_PackedString { t%v, %v }; })", tHashName, value, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, fSizeVariable)
				}
				if parser.UseJs() {
					return fmt.Sprintf("({ char* t%v = static_cast<char*>(static_cast<void*>(%v)); %v })", tHashName, name, cppOutputPackingStringForJs("t"+tHashName, fSizeVariable))
				}
				return fmt.Sprintf("({ char* t%v = static_cast<char*>(static_cast<void*>(%v)); %v_PackedString { t%v, %v }; })", tHashName, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, fSizeVariable)
			}

			if strings.Contains(vOld, "const") {
				if parser.UseJs() {
					return fmt.Sprintf("({ %v pret%v = %v; char* t%v = static_cast<char*>(static_cast<void*>(const_cast<%v*>(&pret%v))); %v })", vOld, tHashName, name, tHashName, value, tHashName, cppOutputPackingStringForJs("t"+tHashName, "-1"))
				}
				return fmt.Sprintf("({ %v pret%v = %v; char* t%v = static_cast<char*>(static_cast<void*>(const_cast<%v*>(&pret%v))); %v_PackedString { t%v, %v }; })", vOld, tHashName, name, tHashName, value, tHashName, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, "-1")
			}
			if parser.UseJs() {
				return fmt.Sprintf("({ %v pret%v = %v; char* t%v = static_cast<char*>(static_cast<void*>(&pret%v)); %v })", vOld, tHashName, name, tHashName, tHashName, cppOutputPackingStringForJs("t"+tHashName, "-1"))
			}
			return fmt.Sprintf("({ %v pret%v = %v; char* t%v = static_cast<char*>(static_cast<void*>(&pret%v)); %v_PackedString { t%v, %v }; })", vOld, tHashName, name, tHashName, tHashName, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, "-1")
		}

	case "QString":
		{
			if strings.Contains(vOld, "*") {
				if parser.UseJs() {
					return fmt.Sprintf("({ QByteArray t%v = %v->toUtf8(); %v })", tHashName, name, cppOutputPackingStringForJs("const_cast<char*>(t"+tHashName+".prepend(\"WHITESPACE\").constData()+10)", "t"+tHashName+".size()-10"))
				}
				return fmt.Sprintf("({ QByteArray t%v = %v->toUtf8(); %v_PackedString { const_cast<char*>(t%v.prepend(\"WHITESPACE\").constData()+10), t%v.size()-10 }; })", tHashName, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, tHashName)
			}
			if parser.UseJs() {
				return fmt.Sprintf("({ QByteArray t%v = %v.toUtf8(); %v })", tHashName, name, cppOutputPackingStringForJs("const_cast<char*>(t"+tHashName+".prepend(\"WHITESPACE\").constData()+10)", "t"+tHashName+".size()-10"))
			}
			return fmt.Sprintf("({ QByteArray t%v = %v.toUtf8(); %v_PackedString { const_cast<char*>(t%v.prepend(\"WHITESPACE\").constData()+10), t%v.size()-10 }; })", tHashName, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, tHashName)
		}

	case "QStringList":
		{
			if strings.Contains(vOld, "*") {
				if parser.UseJs() {
					return fmt.Sprintf("({ QByteArray t%v = %v->join(\"|\").toUtf8(); %v })", tHashName, name, cppOutputPackingStringForJs("const_cast<char*>(t"+tHashName+".prepend(\"WHITESPACE\").constData()+10)", "t"+tHashName+".size()-10"))
				}
				return fmt.Sprintf("({ QByteArray t%v = %v->join(\"|\").toUtf8(); %v_PackedString { const_cast<char*>(t%v.prepend(\"WHITESPACE\").constData()+10), t%v.size()-10 }; })", tHashName, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, tHashName)
			}
			if parser.UseJs() {
				return fmt.Sprintf("({ QByteArray t%v = %v.join(\"|\").toUtf8(); %v })", tHashName, name, cppOutputPackingStringForJs("const_cast<char*>(t"+tHashName+".prepend(\"WHITESPACE\").constData()+10)", "t"+tHashName+".size()-10"))
			}
			return fmt.Sprintf("({ QByteArray t%v = %v.join(\"|\").toUtf8(); %v_PackedString { const_cast<char*>(t%v.prepend(\"WHITESPACE\").constData()+10), t%v.size()-10 }; })", tHashName, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module), tHashName, tHashName)
		}

	case
		"bool", "GLboolean",

		"short", "qint16", "GLshort",
		"ushort", "unsigned short", "quint16", "GLushort",

		"int", "qint32", "GLint", "GLsizei", "GLintptrARB", "GLsizeiptrARB", "GLfixed", "GLclampx",
		"uint", "unsigned int", "quint32", "GLenum", "GLbitfield", "GLuint", "QRgb",

		"long",
		"ulong", "unsigned long",

		"longlong", "long long", "qlonglong", "qint64",
		"ulonglong", "unsigned long long", "qulonglong", "quint64",

		"float", "GLfloat", "GLclampf",
		"double", "qreal",

		"uintptr_t", "uintptr", "quintptr", "WId":
		{
			if strings.Contains(vOld, "*") {
				if value == "bool" || value == "GLboolean" {
					if parser.UseJs() {
						if f.SignalMode == parser.CALLBACK {
							return fmt.Sprintf("reinterpret_cast<uintptr_t>(%v)", name)
						}
						for _, p := range append(f.Parameters, &parser.Parameter{Value: f.Output}) {
							if parser.IsPackedList(p.Value) || parser.IsPackedMap(p.Value) {
								return fmt.Sprintf("reinterpret_cast<uintptr_t>(%v)", name)
							}
							switch parser.CleanValue(p.Value) {
							case "char", "qint8", "uchar", "quint8", "GLubyte", "QString", "QStringList":
								return fmt.Sprintf("reinterpret_cast<uintptr_t>(%v)", name)
							}
						}
					}
					return fmt.Sprintf("reinterpret_cast<char*>(%v)", name)
				}
				return fmt.Sprintf("*%v", name)
			}

			return name
		}

		//non std types

	case "void", "GLvoid", "", "T", "JavaVM", "jclass", "jobject":
		{
			if value == "void" || value == "T" {
				if strings.Contains(vOld, "*") && strings.Contains(vOld, "const") {
					return fmt.Sprintf("const_cast<void*>(%v)", name)
				}
			}

			return name
		}
	}

	switch {
	case isEnum(f.ClassName(), value):
		{
			if parser.UseJs() {
				return fmt.Sprintf("enum_cast<long>(%v)", name)
			}
			return name
		}

	case isClass(value):
		{
			if strings.Contains(value, ".") {
				value = strings.Split(value, ".")[1]
			}
			if strings.Contains(vOld, "*") {
				if strings.Contains(vOld, "const") {
					return fmt.Sprintf("const_cast<%v*>(%v)", value, name)
				}
				return name
			}

			if strings.Contains(vOld, "&") {
				if strings.Contains(vOld, "const") {
					return fmt.Sprintf("const_cast<%v*>(&%v)", value, name)
				}
				if f.SignalMode == parser.CALLBACK {
					return fmt.Sprintf("static_cast<%v*>(&%v)", value, name)
				}
			}

			f.NeedsFinalizer = true

			switch value {
			case "QModelIndex", "QMetaMethod", "QItemSelection":
				{
					return fmt.Sprintf("new %v(%v)", value, name)
				}

			case "QAndroidJniObject":
				{
					return fmt.Sprintf("new %v(%v.object())", value, name)
				}

			case "QPoint", "QPointF":
				{
					return fmt.Sprintf("({ %v tmpValue = %v; new %v(tmpValue.x(), tmpValue.y()); })", value, name, value)
				}

			case "QSize", "QSizeF":
				{
					return fmt.Sprintf("({ %v tmpValue = %v; new %v(tmpValue.width(), tmpValue.height()); })", value, name, value)
				}

			case "QRect", "QRectF":
				{
					return fmt.Sprintf("({ %v tmpValue = %v; new %v(tmpValue.x(), tmpValue.y(), tmpValue.width(), tmpValue.height()); })", value, name, value)
				}

			case "QLine", "QLineF":
				{
					return fmt.Sprintf("({ %v tmpValue = %v; new %v(tmpValue.p1(), tmpValue.p2()); })", value, name, value)
				}

			case "QMargins", "QMarginsF":
				{
					return fmt.Sprintf("({ %v tmpValue = %v; new %v(tmpValue.left(), tmpValue.top(), tmpValue.right(), tmpValue.bottom()); })", value, name, value)
				}
			}

			switch f.Fullname {
			case "QColor::toVariant", "QFont::toVariant", "QImage::toVariant", "QObject::toVariant", "QIcon::toVariant", "QBrush::toVariant":
				{
					if f.Fullname == "QObject::toVariant" {
						return fmt.Sprintf("new %v(QVariant::fromValue(%v))", value, strings.Split(name, "->")[0])
					}
					return fmt.Sprintf("new %v(*%v)", value, strings.Split(name, "->")[0])
				}

			case "QVariant::toColor", "QVariant::toFont", "QVariant::toImage", "QVariant::toObject", "QVariant::toIcon", "QVariant::toBrush":
				{
					f.NeedsFinalizer = false

					if f.Fullname == "QVariant::toObject" {
						return fmt.Sprintf("qvariant_cast<%v*>(*%v)", value, strings.Split(name, "->")[0])
					}
					return fmt.Sprintf("new %v(qvariant_cast<%v>(*%v))", value, value, strings.Split(name, "->")[0])
				}
			}

			for _, f := range parser.State.ClassMap[value].Functions {
				if f.Meta == parser.CONSTRUCTOR {
					switch len(f.Parameters) {
					case 0:
						{
							if value == "QDataStream" {

							} else {
								return fmt.Sprintf("new %v(%v)", value, name)
							}
						}

					case 1:
						{
							if parser.CleanValue(f.Parameters[0].Value) == value {
								return fmt.Sprintf("new %v(%v)", value, name)
							}
						}
					}
				}
			}
		}

	case parser.IsPackedList(value) || parser.IsPackedMap(value):
		{
			if strings.HasSuffix(vOld, "*") {
				if strings.Contains(vOld, "const") {
					if parser.UseJs() {
						return fmt.Sprintf("({ %v* tmpValue = const_cast<%v*>(%v); %v ret; })", value, value, name, cppOutputPackingListForJs())
					}
					return fmt.Sprintf("({ %v* tmpValue = const_cast<%v*>(%v); %v_PackedList { tmpValue, tmpValue->size() }; })", value, value, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module))
				}
				if parser.UseJs() {
					return fmt.Sprintf("({ %v* tmpValue = %v; %v ret; })", value, name, cppOutputPackingListForJs())
				}
				return fmt.Sprintf("({ %v* tmpValue = %v; %v_PackedList { tmpValue, tmpValue->size() }; })", value, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module))
			}

			if strings.HasSuffix(vOld, "&") {
				if strings.Contains(vOld, "const") {
					if parser.UseJs() {
						return fmt.Sprintf("({ %v* tmpValue = const_cast<%v*>(&%v); %v ret; })", value, value, name, cppOutputPackingListForJs())
					}
					return fmt.Sprintf("({ %v* tmpValue = const_cast<%v*>(&%v); %v_PackedList { tmpValue, tmpValue->size() }; })", value, value, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module))
				}
				if f.SignalMode == parser.CALLBACK {
					if parser.UseJs() {
						return fmt.Sprintf("({ %v* tmpValue = static_cast<%v*>(&%v); %v ret; })", value, value, name, cppOutputPackingListForJs())
					}
					return fmt.Sprintf("({ %v* tmpValue = static_cast<%v*>(&%v); %v_PackedList { tmpValue, tmpValue->size() }; })", value, value, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module))
				}
			}

			if parser.UseJs() {
				return fmt.Sprintf("({ %v* tmpValue = new %v(%v); %v ret; })", value, value, name, cppOutputPackingListForJs())
			}
			return fmt.Sprintf("({ %v* tmpValue = new %v(%v); %v_PackedList { tmpValue, tmpValue->size() }; })", value, value, name, strings.Title(parser.State.ClassMap[f.ClassName()].Module))
		}
	}

	f.Access = fmt.Sprintf("unsupported_cppOutput(%v)", value)
	return f.Access
}

func GoOutputJS(name, value string, f *parser.Function, p string) string {
	return goOutputJS(name, value, f, p)
}

func goOutputJS(name, value string, f *parser.Function, p string) string {

	var vOld = value

	name = parser.CleanName(name, value)
	value = parser.CleanValue(value)

	switch value {
	case "char", "qint8", "uchar", "quint8", "GLubyte", "QString":
		{
			return func() string {
				var out = fmt.Sprintf("jsGoUnpackString(%v.Get(\"data\").String())", name)
				if strings.Contains(p, "error") {
					return fmt.Sprintf("errors.New(%v)", out)
				}
				return out
			}()
		}

	case "QStringList":
		{
			if f.SignalMode == parser.CALLBACK {
				return fmt.Sprintf("jsGoUnpackString(%v.Get(\"data\").String())", name)
			}
			return fmt.Sprintf("strings.Split(jsGoUnpackString(%v.Get(\"data\").String()), \"|\")", name)
		}

	case "void", "GLvoid", "":
		{
			if strings.Contains(vOld, "*") {
				if parser.UseWasm() {
					return "unsafe.Pointer(uintptr(" + name + ".Int()))"
				}
				return "unsafe.Pointer(" + name + ")"
			}
			return name
		}

	case "bool", "GLboolean":
		{
			if parser.UseWasm() && f.SignalMode != parser.CALLBACK { //callback arguments for wasm are proper bools, this would panic otherwise: https://github.com/golang/go/blob/master/src/syscall/js/js.go#L361
				return fmt.Sprintf("int8(%v.Int()) != 0", name)
			}
			return fmt.Sprintf("%v.Bool()", name)
		}

	case "short", "qint16", "GLshort":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("int16(%v.Int())", name)
			}
			return fmt.Sprintf("int16(%v.Int64())", name)
		}

	case "ushort", "unsigned short", "quint16", "GLushort":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("uint16(%v.Int())", name)
			}
			return fmt.Sprintf("uint16(%v.Uint64())", name)
		}

	case "int", "qint32", "GLint", "GLsizei", "GLintptrARB", "GLsizeiptrARB", "GLfixed", "GLclampx":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("int(int32(%v.Int()))", name)
			}
			return fmt.Sprintf("int(int32(%v.Int64()))", name)
		}

	case "uint", "unsigned int", "quint32", "GLenum", "GLbitfield", "GLuint", "QRgb":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("uint(uint32(%v.Int()))", name)
			}
			return fmt.Sprintf("uint(uint32(%v.Uint64()))", name)
		}

	case "long":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("int(int32(%v.Int()))", name)
			}
			return fmt.Sprintf("int(int32(%v.Int64()))", name)
		}

	case "ulong", "unsigned long":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("uint(uint32(%v.Int()))", name)
			}
			return fmt.Sprintf("uint(uint32(%v.Uint64()))", name)
		}

	case "longlong", "long long", "qlonglong", "qint64":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("int64(%v.Int())", name)
			}
			return fmt.Sprintf("int64(%v.Int64())", name)
		}

	case "ulonglong", "unsigned long long", "qulonglong", "quint64":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("uint64(%v.Int())", name)
			}
			return fmt.Sprintf("uint64(%v.Uint64())", name)
		}

	case "float", "GLfloat", "GLclampf":
		{
			return fmt.Sprintf("float32(%v.Float())", name)
		}

	case "double", "qreal":
		{
			return fmt.Sprintf("float64(%v.Float())", name)
		}

	case "uintptr_t", "uintptr", "quintptr", "WId":
		{
			if parser.UseJs() {
				if parser.UseWasm() {
					return fmt.Sprintf("uintptr(%v.Int())", name)
				}
				return fmt.Sprintf("uintptr(%v.Unsafe())", name)
			}
			return fmt.Sprintf("uintptr(%v)", name)
		}

		//non std types

	case "T", "JavaVM", "jclass", "jobject":
		{
			switch f.TemplateModeJNI {
			case "Boolean":
				{
					if parser.UseWasm() {
						return fmt.Sprintf("int8(%v.Int()) != 0", name)
					}
					return fmt.Sprintf("%v.Bool()", name)
				}

			case "Int":
				{
					return fmt.Sprintf("int(int32(%v.Uint64()))", name)
				}

			case "Void":
				{
					return name
				}
			}

			return fmt.Sprintf("unsafe.Pointer(%v)", name)
		}
	}

	switch {
	case isEnum(f.ClassName(), value):
		{
			if c, ok := parser.State.ClassMap[class(cppEnum(f, value, false))]; ok && module(c.Module) != module(f) && module(c.Module) != "" {
				if _, ok := parser.State.ClassMap[f.ClassName()].WeakLink[c.Module]; ok {
					if parser.UseWasm() {
						return fmt.Sprintf("int64(%v.Int())", name)
					}
					return fmt.Sprintf("int64(%v.Int64())", name)
				}
				if parser.UseWasm() {
					return fmt.Sprintf("%v.%v(%v.Int())", module(c.Module), goEnum(f, value), name)
				}
				return fmt.Sprintf("%v.%v(%v.Int64())", module(c.Module), goEnum(f, value), name)
			}
			if parser.UseWasm() {
				return fmt.Sprintf("%v(%v.Int())", goEnum(f, value), name)
			}
			return fmt.Sprintf("%v(%v.Int64())", goEnum(f, value), name)
		}

	case isClass(value):
		{
			if parser.UseWasm() && f.SignalMode == parser.CALLBACK {
				return fmt.Sprintf("uintptr(%v.Int())", name)
			}

			if strings.Contains(value, ".") {
				value = strings.Split(value, ".")[1]
			}
			if m := module(parser.State.ClassMap[value].Module); m != module(f) {
				if _, ok := parser.State.ClassMap[f.ClassName()].WeakLink[parser.State.ClassMap[value].Module]; ok {
					if parser.UseWasm() {
						return fmt.Sprintf("unsafe.Pointer(uintptr(%v.Int()))", name)
					}
					return fmt.Sprintf("unsafe.Pointer(%v.Unsafe())", name)
				}
				if parser.UseWasm() {
					return fmt.Sprintf("%v.New%vFromPointer(unsafe.Pointer(uintptr(%v.Int())))", m, strings.Title(value), name)
				}
				return fmt.Sprintf("%v.New%vFromPointer(unsafe.Pointer(%v.Unsafe()))", m, strings.Title(value), name)
			}
			if parser.UseWasm() {
				return fmt.Sprintf("New%vFromPointer(unsafe.Pointer(uintptr(%v.Int())))", strings.Title(value), name)
			}
			return fmt.Sprintf("New%vFromPointer(unsafe.Pointer(%v.Unsafe()))", strings.Title(value), name)
		}

	case parser.IsPackedList(value):
		{
			if parser.UseWasm() {
				return fmt.Sprintf("func(l js.Value)%v{out := make(%v, int(l.Get(\"len\").Int()))\ntmpList := New%vFromPointer(unsafe.Pointer(uintptr(l.Get(\"data\").Int())))\nfor i:=0;i<len(out);i++{ out[i] = tmpList.__%v_atList%v(i) }\nreturn out}(%v)", goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, f.OverloadNumber, name)
			}
			return fmt.Sprintf("func(l *js.Object)%v{out := make(%v, int(l.Get(\"len\").Int()))\ntmpList := New%vFromPointer(unsafe.Pointer(l.Get(\"data\").Unsafe()))\nfor i:=0;i<len(out);i++{ out[i] = tmpList.__%v_atList%v(i) }\nreturn out}(%v)", goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, f.OverloadNumber, name)
		}

	case parser.IsPackedMap(value):
		{
			if parser.UseWasm() {
				return fmt.Sprintf("func(l js.Value)%v{out := make(%v, int(l.Get(\"len\").Int()))\ntmpList := New%vFromPointer(unsafe.Pointer(uintptr(l.Get(\"data\").Int())))\nfor i,v:=range tmpList.__%v_keyList(){ out[v] = tmpList.__%v_atList%v(v, i) }\nreturn out}(%v)", goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, f.Name, f.OverloadNumber, name)
			}
			return fmt.Sprintf("func(l *js.Object)%v{out := make(%v, int(l.Get(\"len\").Int()))\ntmpList := New%vFromPointer(unsafe.Pointer(l.Get(\"data\").Unsafe()))\nfor i,v:=range tmpList.__%v_keyList(){ out[v] = tmpList.__%v_atList%v(v, i) }\nreturn out}(%v)", goType(f, value, p), goType(f, value, p), strings.Title(f.ClassName()), f.Name, f.Name, f.OverloadNumber, name)
		}
	}

	f.Access = fmt.Sprintf("unsupported_goOutput(%v)", value)
	return f.Access
}
