all: webgl_1st

webgl_1st: webgl_server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_1st
	cd examples; ./webgl_test_server

webgl_2d: webgl_server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_2d
	cd examples; ./webgl_test_server

webgl_2di: webgl_server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_2di
	cd examples; ./webgl_test_server

webgl_3d: webgl_server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_3d
	cd examples; ./webgl_test_server

webgl_3di: webgl_server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_3di
	cd examples; ./webgl_test_server

webgl_globe: webgl_server
	cd examples; GOOS=js GOARCH=wasm go build -o ./webgl_test.wasm ./webgl_globe
	cd examples; ./webgl_test_server

webgl_server: 
	cd examples; go build -o ./webgl_test_server ./webgl_test_server.go

opengl_1st: 
	go run ./examples/opengl_1st/main.go

opengl_2d: 
	go run ./examples/opengl_2d/main.go

opengl_2di: 
	go run ./examples/opengl_2di/main.go

opengl_3d: 
	go run ./examples/opengl_3d/main.go

opengl_globe: 
	go run ./examples/opengl_globe/main.go

clean:
	rm ./webgl_test.wasm ./webgl_test_server
