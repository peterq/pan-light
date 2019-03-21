package parser

func (m *Module) remove() {
	m.removeClasses()
}

func (m *Module) removeClasses() {
	for _, c := range SortedClassesForModule(m.Project, false) {
		switch {
		case
			!(c.Access == "public" || c.Access == "protected"),
			c.Name == "qoutputrange":
			delete(State.ClassMap, c.Name)
		}
	}
}
