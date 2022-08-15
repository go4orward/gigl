all: webgl_globe run_webgl

webgl: webgl_1st webgl_2d webgl_2di webgl_3d webgl_3di webgl_globe webgl_server

webgl_1st: 
	cd tutorial/webgl_1st;    GOOS=js GOARCH=wasm go build -o ../webgl_server/webgl_test.wasm .

webgl_2d: 
	cd tutorial/webgl_2d;     GOOS=js GOARCH=wasm go build -o ../webgl_server/webgl_test.wasm .

webgl_2di: 
	cd tutorial/webgl_2di;    GOOS=js GOARCH=wasm go build -o ../webgl_server/webgl_test.wasm .

webgl_3d: 
	cd tutorial/webgl_3d;     GOOS=js GOARCH=wasm go build -o ../webgl_server/webgl_test.wasm .

webgl_3di: 
	cd tutorial/webgl_3di;    GOOS=js GOARCH=wasm go build -o ../webgl_server/webgl_test.wasm .

webgl_globe: 
	cd tutorial/webgl_globe;  GOOS=js GOARCH=wasm go build -o ../webgl_server/webgl_test.wasm .

webgl_server: 
	cd tutorial/webgl_server; go build -o ./webgl_test_server ./webgl_test_server.go

webgl_run: webgl_server
	cd tutorial/webgl_server; ./webgl_test_server


opengl_1st: 
	go run ./tutorial/opengl_1st/opengl_1st.go

opengl_2d: 
	go run ./tutorial/opengl_2d/opengl_2d.go

opengl_2di: 
	go run ./tutorial/opengl_2di/opengl_2di.go

opengl_3d: 
	go run ./tutorial/opengl_3d/opengl_3d.go

opengl_globe: 
	go run ./tutorial/opengl_globe/opengl_globe.go

geometry_viewer:
	mkdir -p build
	cd tutorial/geometry_viewer; go build -o ../../build/geometry_viewer .

clean:
	rm tutorial/webgl_server/webgl_test.wasm tutorial/webgl_server/webgl_test_server
