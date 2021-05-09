all: server 1st

1st: server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_1st
	cd examples; ./webgl_test_server

2d: server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_2d
	cd examples; ./webgl_test_server

2dui: server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_2dui
	cd examples; ./webgl_test_server

3d: server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_3d
	cd examples; ./webgl_test_server

globe: server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_globe
	cd examples; ./webgl_test_server

server: 
	cd examples; go build -o ./webgl_test_server ./webgl_test_server.go

clean:
	rm ./webgl_test.wasm ./webgl_test_server
