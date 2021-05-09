all: server 1st

1st: server
	GOOS=js GOARCH=wasm go build -o examples/webgl_test.wasm examples/webgl1st_example.go
	cd examples; ./webgl_test_server

2d: server
	GOOS=js GOARCH=wasm go build -o examples/webgl_test.wasm examples/webgl2d_example.go
	cd examples; ./webgl_test_server

2dui: server
	GOOS=js GOARCH=wasm go build -o examples/webgl_test.wasm examples/webgl2dui_example.go
	cd examples; ./webgl_test_server

3d: server
	GOOS=js GOARCH=wasm go build -o examples/webgl_test.wasm examples/webgl3d_example.go
	cd examples; ./webgl_test_server

globe: server
	GOOS=js GOARCH=wasm go build -o examples/webgl_test.wasm examples/webglglobe_example.go
	cd examples; ./webgl_test_server

server: 
	go build -o examples/webgl_test_server examples/webgl_test_server.go

clean:
	rm examples/webgl_test.wasm examples/webgl_test_server
