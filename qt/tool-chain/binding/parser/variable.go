package parser

import (
	"fmt"
	"strings"
)

type Variable struct {
	Name     string     `xml:"name,attr"`
	Fullname string     `xml:"fullname,attr"`
	Href     string     `xml:"href,attr"`
	Status   string     `xml:"status,attr"`
	Access   string     `xml:"access,attr"`
	Filepath string     `xml:"filepath,attr"`
	Static   bool       `xml:"static,attr"`
	Output   string     `xml:"type,attr"`
	Brief    string     `xml:"brief,attr"`
	Getter   []struct{} `xml:"getter"`
	Setter   []struct{} `xml:"setter"`

	IsMocSynthetic bool
	PureGoType     string
	Connect        int
	ConnectGet     bool
	ConnectSet     bool
	ConnectChanged bool
	Target         string
	Inbound        bool
}

func (v *Variable) Class() (*Class, bool) {
	var class, ok = State.ClassMap[v.ClassName()]
	return class, ok
}

func (v *Variable) ClassName() string {
	var s = strings.Split(v.Fullname, "::")
	if len(s) == 3 {
		return s[1]
	}
	return s[0]
}

func (v *Variable) varToFunc() []*Function {
	var funcs = make([]*Function, 0)

	var class, ok = v.Class()
	if !ok || class.HasFunctionWithName(v.Name) {
		return funcs
	}

	funcs = append(funcs, &Function{
		Name:     v.Name,
		Fullname: v.Fullname,
		Href:     v.Href,
		Status:   v.Status,
		Access:   v.Access,
		Filepath: v.Filepath,
		Static:   v.Static,
		Output:   v.Output,
		Meta:     GETTER,
		Brief:    v.Brief,
	})

	if strings.Contains(v.Output, "const") {
		return funcs
	}

	funcs = append(funcs, &Function{
		Name:       fmt.Sprintf("set%v", strings.Title(v.Name)),
		Fullname:   fmt.Sprintf("%v::set%v", v.ClassName(), strings.Title(v.Name)),
		Href:       v.Href,
		Status:     v.Status,
		Access:     v.Access,
		Filepath:   v.Filepath,
		Static:     v.Static,
		Output:     "void",
		Meta:       SETTER,
		Parameters: []*Parameter{{Value: v.Output}},
		TmpName:    v.Name,
		Brief:      v.Brief,
	})

	return funcs
}

func (v *Variable) propToFunc(c *Class) []*Function {
	var funcs = make([]*Function, 0)

	if len(v.Getter) != 0 {
		return funcs
	}

	if !(c.HasFunctionWithName(v.Name) ||
		c.HasFunctionWithName(fmt.Sprintf("is%v", strings.Title(v.Name))) ||
		c.HasFunctionWithName(fmt.Sprintf("has%v", strings.Title(v.Name)))) {

		tmpF := &Function{
			Name:     v.Name,
			Fullname: v.Fullname,
			Href:     v.Href,
			Status:   v.Status,
			Access:   v.Access,
			Filepath: v.Filepath,
			Static:   v.Static,
			Output:   v.Output,
			Meta:     PLAIN,
			Virtual: func() string {
				if c.Module == MOC && !v.IsMocSynthetic {
					return IMPURE
				}
				return ""
			}(),
			Signature:     "()",
			IsMocFunction: c.Module == MOC,
			IsMocProperty: c.Module == MOC,
		}

		if tmpF.Output == "bool" {
			if !strings.HasPrefix(strings.ToLower(v.Name), "is") {
				tmpF.Name = fmt.Sprintf("is%v", strings.Title(tmpF.Name))
			}
			tmpF.Fullname = fmt.Sprintf("%v::%v", tmpF.ClassName(), tmpF.Name)
		}

		funcs = append(funcs, tmpF)
	}

	if len(v.Setter) != 0 || c.HasFunctionWithName(fmt.Sprintf("set%v", strings.Title(v.Name))) {
		return funcs
	}

	funcs = append(funcs, &Function{
		Name:     fmt.Sprintf("set%v", strings.Title(v.Name)),
		Fullname: fmt.Sprintf("%v::set%v", v.ClassName(), strings.Title(v.Name)),
		Href:     v.Href,
		Status:   v.Status,
		Access:   v.Access,
		Filepath: v.Filepath,
		Static:   v.Static,
		Output:   "void",
		Meta:     PLAIN,
		Virtual: func() string {
			if c.Module == MOC && !v.IsMocSynthetic {
				return IMPURE
			}
			return ""
		}(),
		Parameters:    []*Parameter{{Name: v.Name, Value: v.Output}},
		Signature:     "()",
		IsMocFunction: c.Module == MOC,
		IsMocProperty: c.Module == MOC,
	})

	if c.Module == MOC {
		funcs = append(funcs, &Function{
			Name:          fmt.Sprintf("%vChanged", v.Name),
			Fullname:      fmt.Sprintf("%v::%vChanged", v.ClassName(), v.Name),
			Status:        v.Status,
			Access:        v.Access,
			Output:        "void",
			Meta:          SIGNAL,
			Parameters:    []*Parameter{{Name: v.Name, Value: v.Output}},
			Signature:     "()",
			IsMocFunction: true,
		})
	}

	//add all overloaded property functions from base classes
	//TODO: move rest into seperate function, as this func is called multiple times

	for _, bc := range c.GetAllBases() {
		var bclass, ok = State.ClassMap[bc]
		if !ok {
			continue
		}

		for _, bcf := range bclass.Functions {
			if bcf.Name != fmt.Sprintf("set%v", strings.Title(v.Name)) || !bcf.Overload {
				continue
			}

			var tmpF = *bcf

			tmpF.Name = fmt.Sprintf("set%v", strings.Title(v.Name))
			tmpF.Fullname = fmt.Sprintf("%v::%v", v.ClassName(), tmpF.Name)

			funcs = append(funcs, &tmpF)
		}
	}

	return funcs
}
