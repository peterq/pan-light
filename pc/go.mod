module github.com/peterq/pan-light/pc

require (
	github.com/peterq/pan-light/qt v0.0.0
	github.com/peterq/pan-light/qt/bindings v0.0.0
)

replace github.com/peterq/pan-light/qt/bindings => ../qt/bindings

replace github.com/peterq/pan-light/qt => ../qt
