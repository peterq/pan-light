module github.com/peterq/pan-light/pc

require (
	github.com/golang/protobuf v1.3.1
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/peterq/pan-light/qt v0.0.0
	github.com/peterq/pan-light/qt/bindings v0.0.0
	github.com/peterq/pan-light/qt/tool-chain/binding/files/docs/5.12.0 v0.0.0-00010101000000-000000000000 // indirect
	github.com/pkg/errors v0.8.1
	github.com/therecipe/qt v0.0.0-20190608010047-be44906910f4 // indirect
)

replace github.com/peterq/pan-light/qt/bindings => ../qt/bindings

replace github.com/peterq/pan-light/qt => ../qt

replace github.com/peterq/pan-light/qt/tool-chain/binding/files/docs/5.12.0 => ../qt/tool-chain/binding/files/docs/5.12.0
