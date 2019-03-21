package converter

//TODO: GLchar, GLbyte

import (
	"fmt"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
)

func GoInput(name, value string, f *parser.Function, p string) string {
	if parser.UseJs() {
		return GoInputJS(name, value, f, p)
	}

	var vOld = value

	name = parser.CleanName(name, value)
	value = parser.CleanValue(value)

	switch value {
	case "char", "qint8":
		{
			if strings.Contains(vOld, "**") {
				return fmt.Sprintf("C.CString(strings.Join(%v, \"|\"))", name)
			}

			if value == "char" && strings.Count(vOld, "*") == 1 && f.Name == "readData" {
				return fmt.Sprintf("C.CString(strings.Repeat(\"0\", int(%v)))", parser.CleanName(f.Parameters[1].Name, f.Parameters[1].Value))
			}

			switch value {
			case "char", "qint8":
				if len(f.Parameters) <= 4 &&
					(strings.Contains(strings.ToLower(f.Name), "read") ||
						strings.Contains(strings.ToLower(f.Name), "write") ||
						strings.Contains(strings.ToLower(f.Name), "data")) {
					for _, p := range f.Parameters {
						if strings.Contains(p.Value, "int") && f.Parameters[0].Value == vOld {
							return fmt.Sprintf("(*C.char)(unsafe.Pointer(&%v[0]))", name)
						}
					}
				}
			}

			return fmt.Sprintf("C.CString(%v)", name)
		}

	case "uchar", "quint8", "GLubyte", "QString":
		{
			switch value {
			case "uchar", "quint8", "GLubyte":
				if len(f.Parameters) <= 4 &&
					(strings.Contains(strings.ToLower(f.Name), "read") ||
						strings.Contains(strings.ToLower(f.Name), "write") ||
						strings.Contains(strings.ToLower(f.Name), "data")) {
					for _, p := range f.Parameters {
						if strings.Contains(p.Value, "int") && f.Parameters[0].Value == vOld {
							return fmt.Sprintf("(*C.char)(unsafe.Pointer(&%v[0]))", name)
						}
					}
				}
			}

			return fmt.Sprintf("C.CString(%v)", func() string {
				if strings.Contains(p, "error") {
					return fmt.Sprintf("func() string { tmp := %v\n if tmp != nil { return tmp.Error() }\n return \"\" }()", name)
				}
				return name
			}())
		}

	case "QStringList":
		{
			return fmt.Sprintf("C.CString(strings.Join(%v, \"|\"))", name)
		}

	case "void", "GLvoid" /*, ""*/ :
		{
			if strings.Contains(vOld, "*") {
				return name
			}
		}

	case "bool", "GLboolean":
		{
			if strings.Contains(vOld, "*") {
				return fmt.Sprintf("C.char(int8(qt.GoBoolToInt(*%v)))", name)
			}
			return fmt.Sprintf("C.char(int8(qt.GoBoolToInt(%v)))", name)
		}

	case "short", "qint16", "GLshort":
		{
			return fmt.Sprintf("C.short(%v)", name)
		}

	case "ushort", "unsigned short", "quint16", "GLushort":
		{
			return fmt.Sprintf("C.ushort(%v)", name)
		}

	case "int", "qint32", "GLint", "GLsizei", "GLintptrARB", "GLsizeiptrARB", "GLfixed", "GLclampx":
		{
			return fmt.Sprintf("C.int(int32(%v))", name)
		}

	case "uint", "unsigned int", "quint32", "GLenum", "GLbitfield", "GLuint", "QRgb":
		{
			return fmt.Sprintf("C.uint(uint32(%v))", name)
		}

	case "long":
		{
			return fmt.Sprintf("C.long(int32(%v))", name)
		}

	case "ulong", "unsigned long":
		{
			return fmt.Sprintf("C.ulong(uint32(%v))", name)
		}

	case "longlong", "long long", "qlonglong", "qint64":
		{
			return fmt.Sprintf("C.longlong(%v)", name)
		}

	case "ulonglong", "unsigned long long", "qulonglong", "quint64":
		{
			return fmt.Sprintf("C.ulonglong(%v)", name)
		}

	case "float", "GLfloat", "GLclampf":
		{
			return fmt.Sprintf("C.float(%v)", name)
		}

	case "double", "qreal":
		{
			if value == "qreal" && strings.HasPrefix(parser.State.Target, "sailfish") {
				return fmt.Sprintf("C.float(%v)", name)
			}
			return fmt.Sprintf("C.double(%v)", name)
		}

	case "uintptr_t", "uintptr", "quintptr", "WId":
		{
			return fmt.Sprintf("C.uintptr_t(%v)", name)
		}

		//non std types

	case "T":
		{
			switch f.TemplateModeJNI {
			case "Boolean":
				{
					return fmt.Sprintf("C.char(int8(qt.GoBoolToInt(%v)))", name)
				}

			case "Int":
				{
					return fmt.Sprintf("C.int(int32(%v))", name)
				}
			}

			if module(f) == "androidextras" {
				return "p0"
			}
		}

	case "JavaVM", "jclass", "jobject":
		{
			return name
		}

	case "...":
		{
			var tmp = make([]string, 10)
			for i := range tmp {
				tmp[i] = fmt.Sprintf("p%v", i)
			}
			return strings.Join(tmp, ", ")
		}
	}

	switch {
	case isEnum(f.ClassName(), value):
		{
			return fmt.Sprintf("C.longlong(%v)", name)
		}

	case isClass(value):
		{
			if strings.Contains(value, ".") {
				value = strings.Split(value, ".")[1]
			}
			if m := module(parser.State.ClassMap[value].Module); m != module(f) {
				if _, ok := parser.State.ClassMap[f.ClassName()].WeakLink[parser.State.ClassMap[value].Module]; ok {
					return name
				}
				return fmt.Sprintf("%v.PointerFrom%v(%v)", m, strings.Title(value), name)
			}
			return fmt.Sprintf("PointerFrom%v(%v)", strings.Title(value), name)
		}

	case parser.IsPackedList(value):
		{
			if strings.ContainsAny(name, "*&()[]") {
				return fmt.Sprintf("func() unsafe.Pointer {\ntmpList := New%vFromPointer(New%vFromPointer(nil).__%v_newList%v())\nfor _,v := range %v{\ntmpList.__%v_setList%v(v)\n}\nreturn tmpList.Pointer()\n}()", strings.Title(f.ClassName()), strings.Title(f.ClassName()), f.Name, f.OverloadNumber, name, f.Name, f.OverloadNumber)
			}
			return fmt.Sprintf("func() unsafe.Pointer {\ntmpList := New%vFromPointer(New%vFromPointer(nil).__%v_%v_newList%v())\nfor _,v := range %v{\ntmpList.__%v_%v_setList%v(v)\n}\nreturn tmpList.Pointer()\n}()", strings.Title(f.ClassName()), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, name, f.Name, name, f.OverloadNumber)
		}

	case parser.IsPackedMap(value):
		{
			if strings.ContainsAny(name, "*&()[]") {
				return fmt.Sprintf("func() unsafe.Pointer {\ntmpList := New%vFromPointer(New%vFromPointer(nil).__%v_newList%v())\nfor k,v := range %v{\ntmpList.__%v_setList%v(k, v)\n}\nreturn tmpList.Pointer()\n}()", strings.Title(f.ClassName()), strings.Title(f.ClassName()), f.Name, f.OverloadNumber, name, f.Name, f.OverloadNumber)
			}
			return fmt.Sprintf("func() unsafe.Pointer {\ntmpList := New%vFromPointer(New%vFromPointer(nil).__%v_%v_newList%v())\nfor k,v := range %v{\ntmpList.__%v_%v_setList%v(k, v)\n}\nreturn tmpList.Pointer()\n}()", strings.Title(f.ClassName()), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, name, f.Name, name, f.OverloadNumber)
		}
	}

	f.Access = fmt.Sprintf("unsupported_goInput(%v)", value)
	return f.Access
}

func CppInput(name, value string, f *parser.Function) string {

	if (f.SignalMode == parser.CALLBACK || strings.HasPrefix(name, "callback") || strings.HasPrefix(name, "emscripten::val::global")) && (parser.CleanValue(value) == "QString" || parser.CleanValue(value) == "QStringList") {
		if parser.UseJs() {
			if parser.UseWasm() {
				return fmt.Sprintf("({ emscripten::val tempVal = %v; %v ret = %v; emscripten::val::global(\"Module\").call<void>(\"_callbackReleaseTypedArray\", tempVal[\"data_ptr\"].as<uintptr_t>()); ret; })", name, parser.CleanValue(value), cppInput("tempVal", value, f))
			}
			return fmt.Sprintf("({ emscripten::val tempVal = %v; %v ret = %v; ret; })", name, parser.CleanValue(value), cppInput("tempVal", value, f))
		}
		return fmt.Sprintf("({ %v_PackedString tempVal = %v; %v ret = %v; free(tempVal.data); ret; })", strings.Title(parser.State.ClassMap[f.ClassName()].Module), name, parser.CleanValue(value), cppInput("tempVal", value, f))
	}

	out := cppInput(name, value, f)

	if parser.UseJs() {
		if isEnum(f.ClassName(), parser.CleanValue(value)) {
			out = strings.Replace(out, "static_cast", "enum_cast", -1)
		}
	}

	return out
}

func cppInput(name, value string, f *parser.Function) string {
	var vOld = value

	name = parser.CleanName(name, value)
	value = parser.CleanValue(value)

	switch value {
	case "char", "qint8":
		{
			if strings.Contains(vOld, "**") && name == "argv" {
				return "argvs"
			}

			if parser.UseJs() {
				if strings.Contains(vOld, "*") {
					if strings.Contains(vOld, "const") {
						return fmt.Sprintf("QByteArray::fromStdString(%v[\"data\"].as<std::string>()).constData()", name)
					}
					return fmt.Sprintf("const_cast<char*>(QByteArray::fromStdString(%v[\"data\"].as<std::string>()).constData())", name)
				}
				return fmt.Sprintf("*const_cast<char*>(QByteArray::fromStdString(%v[\"data\"].as<std::string>()).constData())", name)
			}
			if strings.Contains(vOld, "*") {
				if strings.Contains(vOld, "const") {
					return fmt.Sprintf("const_cast<const %v*>(%v)", value, name)
				}
				return name
			}

			return fmt.Sprintf("*%v", name)
		}

	case "uchar", "quint8", "GLubyte":
		{
			if parser.UseJs() {
				if strings.Contains(vOld, "*") {
					if strings.Contains(vOld, "const") {
						return fmt.Sprintf("const_cast<const %v*>(static_cast<%v*>(static_cast<void*>(const_cast<char*>(QByteArray::fromStdString(%v[\"data\"].as<std::string>()).constData()))))", value, value, name)
					}
					return fmt.Sprintf("static_cast<%v*>(static_cast<void*>(const_cast<char*>(QByteArray::fromStdString(%v[\"data\"].as<std::string>()).constData())))", value, name)
				}
				return fmt.Sprintf("*static_cast<%v*>(static_cast<void*>(const_cast<char*>(QByteArray::fromStdString(%v[\"data\"].as<std::string>()).constData())))", value, name)
			}
			if strings.Contains(vOld, "*") {
				if strings.Contains(vOld, "const") {
					return fmt.Sprintf("const_cast<const %v*>(static_cast<%v*>(static_cast<void*>(%v)))", value, value, name)
				}
				return fmt.Sprintf("static_cast<%v*>(static_cast<void*>(%v))", value, name)
			}
			return fmt.Sprintf("*static_cast<%v*>(static_cast<void*>(%v))", value, name)
		}

	case "QString":
		{
			if strings.Contains(vOld, "*") {
				if parser.UseJs() {
					return fmt.Sprintf("new QString(QString::fromStdString(%v[\"data\"].as<std::string>()))", name)
				}
				return fmt.Sprintf("new QString(QString::fromUtf8(%[1]v.data, %[1]v.len))", name)
			}

			if strings.Contains(vOld, "&") && !strings.Contains(vOld, "const") {
				return fmt.Sprintf("*(%v)", cppInput(name, "QString*", f))
			}

			if parser.UseJs() {
				return fmt.Sprintf("QString::fromStdString(%v[\"data\"].as<std::string>())", name)
			}
			return fmt.Sprintf("QString::fromUtf8(%[1]v.data, %[1]v.len)", name)
		}

	case "QStringList":
		{
			if strings.Contains(vOld, "*") {
				return fmt.Sprintf("new QStringList(%v)", cppInput(name, "QStringList", f))
			}

			if strings.Contains(vOld, "&") && !strings.Contains(vOld, "const") {
				return fmt.Sprintf("*(%v)", cppInput(name, "QStringList*", f))
			}

			if parser.UseJs() {
				return fmt.Sprintf("QString::fromStdString(%v[\"data\"].as<std::string>()).split(\"|\", QString::SkipEmptyParts)", name)
			}
			return fmt.Sprintf("QString::fromUtf8(%[1]v.data, %[1]v.len).split(\"|\", QString::SkipEmptyParts)", name)
		}

	case "void", "GLvoid" /*, ""*/ :
		{
			if strings.Count(vOld, "*") == 2 && !strings.Contains(vOld, "**") {
				break
			}

			if strings.Contains(vOld, "**") {
				return fmt.Sprintf("&%v", name)
			}

			if parser.UseJs() {
				if strings.Contains(vOld, "*") {
					return fmt.Sprintf("reinterpret_cast<void*>(%v)", name)
				}
			}

			if strings.Contains(vOld, "*") {
				return name
			}
		}

	case "bool", "GLboolean":
		{
			if strings.Contains(vOld, "*") {
				if parser.UseJs() && f.SignalMode == parser.CALLBACK {
					return fmt.Sprintf("reinterpret_cast<uintptr_t>(%v)", value, name)
				}
				return fmt.Sprintf("reinterpret_cast<%v*>(%v)", value, name)
			}
			return fmt.Sprintf("%v != 0", name)
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
			if parser.UseJs() && f.BoundByEmscripten {
				switch value {
				case "longlong", "long long", "qlonglong", "qint64",
					"ulonglong", "unsigned long long", "qulonglong", "quint64":
					f.BoundByEmscripten = false
					name = fmt.Sprintf("enum_cast<%v>(%v)", cppType(f, value), name)
					f.BoundByEmscripten = true
				}
			}
			if strings.Contains(vOld, "&") && name == "argc" {
				return "argcs"
			}

			if strings.Contains(vOld, "*") {
				if strings.Contains(vOld, "const") {
					return fmt.Sprintf("const_cast<const %v*>(&%v)", value, name)
				}
				return fmt.Sprintf("&%v", name)
			}

			return name
		}

		//non std types

	case "T":
		{
			switch f.TemplateModeJNI {
			case "Boolean", "Int":
				{
					return name
				}
			}

			if module(f) == "androidextras" {
				return fmt.Sprintf("static_cast<jobject>(%v)", name)
			}
		}

	case "JavaVM", "jclass", "jobject":
		{
			return fmt.Sprintf("static_cast<%v>(%v)", value, name)
		}

	case "...":
		{
			var tmp = make([]string, 10)
			for i := range tmp {
				if i == 9 {
					tmp[i] = fmt.Sprintf("static_cast<jobject>(%v)", name)
				} else {
					tmp[i] = fmt.Sprintf("static_cast<jobject>(%v%v)", name, i)
				}
			}
			return strings.Join(tmp, ", ")
		}
	}

	switch {
	case isEnum(f.ClassName(), value):
		{
			if !strings.Contains(vOld, "*") {
				if f.Meta == parser.SLOT && f.SignalMode == "" && value == "Qt::Alignment" {
					return fmt.Sprintf("static_cast<%v>(static_cast<%v>(%v))", value, cppEnum(f, value, false), name)
				}
				return fmt.Sprintf("static_cast<%v>(%v)", cppEnum(f, value, false), name)
			}
		}

	case isClass(value):
		{
			if strings.Contains(value, ".") {
				value = strings.Split(value, ".")[1]
			}
			if strings.Contains(vOld, "*") && strings.Contains(vOld, "&") {
				break
			}

			if parser.State.ClassMap[value].Fullname != "" {
				value = parser.State.ClassMap[value].Fullname
			}

			if strings.Contains(vOld, "*") {
				return fmt.Sprintf("static_cast<%v*>(%v)", value, name)
			}
			return fmt.Sprintf("*static_cast<%v*>(%v)", value, name)
		}

	case parser.IsPackedList(value) || parser.IsPackedMap(value):
		{
			if strings.HasSuffix(vOld, "*") {
				return fmt.Sprintf("static_cast<%v*>(%v)", value, name)
			}

			if strings.HasPrefix(vOld, "const") || f.Fullname == "QMacToolBar::setItems" || f.Fullname == "QMacToolBar::setAllowedItems" {
				return fmt.Sprintf("*static_cast<%v*>(%v)", value, name)
			}

			return fmt.Sprintf("({ %v* tmpP = static_cast<%v*>(%v); %v tmpV = *tmpP; tmpP->~%v(); free(tmpP); tmpV; })", parser.CleanValue(value), value, name, parser.CleanValue(value), strings.Split(parser.CleanValue(value), "<")[0])
		}
	}

	f.Access = fmt.Sprintf("unsupported_cppInput(%v)", value)
	return f.Access
}

func GoInputJS(name, value string, f *parser.Function, p string) string {
	var vOld = value

	name = parser.CleanName(name, value)
	value = parser.CleanValue(value)

	switch value {
	case "char", "qint8", "uchar", "quint8", "GLubyte", "QString":
		{
			if strings.Contains(p, "error") {
				name = fmt.Sprintf("func() string {\ntmp := %v\nif tmp != nil { return tmp.Error() }\nreturn \"\"\n}()", name)
			}

			if strings.Contains(vOld, "**") {
				if parser.UseWasm() {
					return fmt.Sprintf("func() js.Value {\ntmp := js.TypedArrayOf([]byte(strings.Join(%v, \"|\")))\nreturn js.ValueOf(map[string]interface{}{\"data\": tmp, \"data_ptr\": unsafe.Pointer(&tmp)})\n}()", name)
				}
				if f.SignalMode != parser.CALLBACK {
					return fmt.Sprintf("func() *js.Object {\ntmp := new(js.Object)\nif js.InternalObject(%v).Get(\"$val\") == js.Undefined {\ntmp.Set(\"data\", []byte(js.InternalObject(%v).Call(\"join\", \"|\").String()))\n} else {\ntmp.Set(\"data\", []byte(strings.Join(%v, \"|\")))\n}\nreturn tmp\n}()", name, name, name) //needed for indirect exported pure js call -> can be ommited if build without js support
				}
				return fmt.Sprintf("func() *js.Object {\ntmp := new(js.Object)\ntmp.Set(\"data\", []byte(strings.Join(%v, \"|\")))\nreturn tmp\n}()", name)
			}

			if value == "char" && strings.Count(vOld, "*") == 1 && f.Name == "readData" {
				//TODO:
			}

			if parser.UseWasm() {
				return fmt.Sprintf("func() js.Value {\ntmp := js.TypedArrayOf([]byte(%v))\nreturn js.ValueOf(map[string]interface{}{\"data\": tmp, \"data_ptr\": unsafe.Pointer(&tmp)})\n}()", name)
			}
			return fmt.Sprintf("func() *js.Object {\ntmp := new(js.Object)\ntmp.Set(\"data\", []byte(%v))\nreturn tmp\n}()", name)
		}

	case "QStringList":
		{
			if parser.UseWasm() {
				return fmt.Sprintf("func() js.Value {\ntmp := js.TypedArrayOf([]byte(strings.Join(%v, \"|\")))\nreturn js.ValueOf(map[string]interface{}{\"data\": tmp, \"data_ptr\": unsafe.Pointer(&tmp)})\n}()", name)
			}
			if f.SignalMode != parser.CALLBACK {
				return fmt.Sprintf("func() *js.Object {\ntmp := new(js.Object)\nif js.InternalObject(%v).Get(\"$val\") == js.Undefined {\ntmp.Set(\"data\", []byte(js.InternalObject(%v).Call(\"join\", \"|\").String()))\n} else {\ntmp.Set(\"data\", []byte(strings.Join(%v, \"|\")))\n}\nreturn tmp\n}()", name, name, name) //needed for indirect exported pure js call -> can be ommited if build without js support
			}
			return fmt.Sprintf("func() *js.Object {\ntmp := new(js.Object)\ntmp.Set(\"data\", []byte(strings.Join(%v, \"|\")))\nreturn tmp\n}()", name)
		}

	case "void", "GLvoid", "":
		{
			if strings.Contains(vOld, "*") {
				return fmt.Sprintf("uintptr(%v)", name)
			}
			return name
		}

	case "bool", "GLboolean":
		{
			return fmt.Sprintf("%v", name)
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
					if parser.UseWasm() || f.SignalMode == parser.CALLBACK {
						return fmt.Sprintf("int64(%v)", name)
					} else {
						return fmt.Sprintf("func() int64 {\nif js.InternalObject(%v).Get(\"$val\") == js.Undefined {\nreturn int64(js.InternalObject(%v).Int64())\n}\nreturn int64(%v)\n}()", name, name, name) //needed for indirect exported pure js call -> can be ommited if build without js support
					}
				}
				if parser.UseWasm() {
					return fmt.Sprintf("int64(%v.%v(%v))", module(c.Module), goEnum(f, value), name)
				}
				if f.SignalMode != parser.CALLBACK {
					return fmt.Sprintf("func() %[1]v.%[2]v {\nif js.InternalObject(%[3]v).Get(\"$val\") == js.Undefined {\nreturn %[1]v.%[2]v(js.InternalObject(%[3]v).Int64())\n}\nreturn %[1]v.%[2]v(%[3]v)\n}()", module(c.Module), goEnum(f, value), name) //needed for indirect exported pure js call -> can be ommited if build without js support
				}
				return fmt.Sprintf("%v.%v(%v)", module(c.Module), goEnum(f, value), name)
			}
			if parser.UseWasm() {
				return fmt.Sprintf("int64(%v(%v))", goEnum(f, value), name)
			}
			if f.SignalMode != parser.CALLBACK {
				return fmt.Sprintf("func() %[1]v {\nif js.InternalObject(%[2]v).Get(\"$val\") == js.Undefined {\nreturn %[1]v(js.InternalObject(%[2]v).Int64())\n}\nreturn %[1]v(%[2]v)\n}()", goEnum(f, value), name) //needed for indirect exported pure js call -> can be ommited if build without js support
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
					return fmt.Sprintf("uintptr(unsafe.Pointer(%v))", name)
				}
				return fmt.Sprintf("uintptr(%v.PointerFrom%v(%v))", m, strings.Title(value), name)
			}
			return fmt.Sprintf("uintptr(PointerFrom%v(%v))", strings.Title(value), name)
		}

	case parser.IsPackedList(value):
		{
			if strings.ContainsAny(name, "*&()[]") {
				return fmt.Sprintf("func() uintptr {\ntmpList := New%vFromPointer(unsafe.Pointer(New%vFromPointer(nil).__%v_newList%v()))\nfor _,v := range %v{\ntmpList.__%v_setList%v(v)\n}\nreturn uintptr(tmpList.Pointer())\n}()", strings.Title(f.ClassName()), strings.Title(f.ClassName()), f.Name, f.OverloadNumber, name, f.Name, f.OverloadNumber)
			}
			return fmt.Sprintf("func() uintptr {\ntmpList := New%vFromPointer(unsafe.Pointer(New%vFromPointer(nil).__%v_%v_newList%v()))\nfor _,v := range %v{\ntmpList.__%v_%v_setList%v(v)\n}\nreturn uintptr(tmpList.Pointer())\n}()", strings.Title(f.ClassName()), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, name, f.Name, name, f.OverloadNumber)
		}

	case parser.IsPackedMap(value):
		{
			if strings.ContainsAny(name, "*&()[]") {
				return fmt.Sprintf("func() uintptr {\ntmpList := New%vFromPointer(unsafe.Pointer(New%vFromPointer(nil).__%v_newList%v()))\nfor k,v := range %v{\ntmpList.__%v_setList%v(k, v)\n}\nreturn uintptr(tmpList.Pointer())\n}()", strings.Title(f.ClassName()), strings.Title(f.ClassName()), f.Name, f.OverloadNumber, name, f.Name, f.OverloadNumber)
			}
			return fmt.Sprintf("func() uintptr {\ntmpList := New%vFromPointer(unsafe.Pointer(New%vFromPointer(nil).__%v_%v_newList%v()))\nfor k,v := range %v{\ntmpList.__%v_%v_setList%v(k, v)\n}\nreturn uintptr(tmpList.Pointer())\n}()", strings.Title(f.ClassName()), strings.Title(f.ClassName()), f.Name, name, f.OverloadNumber, name, f.Name, name, f.OverloadNumber)
		}
	}

	f.Access = fmt.Sprintf("unsupported_goInputJS(%v)", value)
	return f.Access
}
