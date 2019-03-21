package parser

import "strings"

type Enum struct {
	Name     string   `xml:"name,attr"`
	Fullname string   `xml:"fullname,attr"`
	Status   string   `xml:"status,attr"`
	Access   string   `xml:"access,attr"`
	Typedef  string   `xml:"typedef,attr"`
	Values   []*Value `xml:"value"`
	NoConst  bool
}

type Value struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (e *Enum) Class() (*Class, bool) {
	var class, ok = State.ClassMap[e.ClassName()]
	return class, ok
}

func (e *Enum) ClassName() string {
	return strings.Split(e.Fullname, "::")[0]
}

func (e *Enum) register(m string) {

	if c, ok := e.Class(); !ok {
		State.ClassMap[e.ClassName()] = &Class{
			Name:   e.ClassName(),
			Status: "commendable",
			Module: m,
			Access: "public",
			Enums:  []*Enum{e},
		}
	} else {
		c.Enums = append(c.Enums, e)
	}
}
