all: server 1st

1st: server
	GOOS=js GOARCH=wasm go build -o ./_examples/webgl_test.wasm ./_examples/webgl_1st
	cd _examples; ./webgl_test_server

2d: server
	GOOS=js GOARCH=wasm go build -o ./_examples/webgl_test.wasm ./_examples/webgl_2d
	cd _examples; ./webgl_test_server

2dui: server
	GOOS=js GOARCH=wasm go build -o ./_examples/webgl_test.wasm ./_examples/webgl_2dui
	cd _examples; ./webgl_test_server

3d: server
	GOOS=js GOARCH=wasm go build -o ./_examples/webgl_test.wasm ./_examples/webgl_3d
	cd _examples; ./webgl_test_server

globe: server
	GOOS=js GOARCH=wasm go build -o ./_examples/webgl_test.wasm ./_examples/webgl_globe
	cd _examples; ./webgl_test_server

server: 
	go build -o ./_examples/webgl_test_server ./_examples/webgl_test_server.go

clean:
	rm ./_examples/webgl_test.wasm ./_examples/webgl_test_server
