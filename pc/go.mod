module github.com/peterq/pan-light/pc

require (
	github.com/golang/protobuf v1.3.1
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/peterq/pan-light/qt v0.0.0
	github.com/peterq/pan-light/qt/bindings v0.0.0
	github.com/pkg/errors v0.8.1
	github.com/therecipe/qt v0.0.0-20190608010047-be44906910f4 // indirect
)

replace github.com/peterq/pan-light/qt/bindings => ../qt/bindings

replace github.com/peterq/pan-light/qt => ../qt
