package main

import (
	"fmt"

	"github.com/go4orward/gigl/env/webgl10"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/webgl_example.go'.
	fmt.Println("Hello WebGL 1.0")                       // printed in the browser console
	wcanvas, err := webgl10.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		fmt.Printf("Failed to start WebGL : %v\n", err)
		return
	}
	rc := wcanvas.GetGLRenderingContext()
	context, c := wcanvas.RcGetWebGLRenderingContext(), rc.GetConstants()

	// Geometry
	vertices := []float32{-0.5, 0.5, 0, -0.5, -0.5, 0, 0.5, -0.5, 0}
	indices := []uint32{2, 1, 0}
	var vertices_array = wcanvas.RcConvertGoSliceToJsTypedArray(vertices)
	var indices_array = wcanvas.RcConvertGoSliceToJsTypedArray(indices)
	vertexBuffer := context.Call("createBuffer", c.ARRAY_BUFFER)                     // create buffer
	context.Call("bindBuffer", c.ARRAY_BUFFER, vertexBuffer)                         // bind the buffer
	context.Call("bufferData", c.ARRAY_BUFFER, vertices_array, c.STATIC_DRAW)        // pass data to buffer
	indexBuffer := context.Call("createBuffer", c.ELEMENT_ARRAY_BUFFER)              // create index buffer
	context.Call("bindBuffer", c.ELEMENT_ARRAY_BUFFER, indexBuffer)                  // bind the buffer
	context.Call("bufferData", c.ELEMENT_ARRAY_BUFFER, indices_array, c.STATIC_DRAW) // pass data to the buffer

	//// Shaders ////
	vshader_source := `
		attribute vec3 xyz;
		void main(void) {
			gl_Position = vec4(xyz, 1.0);
		}`
	fshader_source := `
		void main(void) {
			gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
		}`
	vshader := context.Call("createShader", c.VERTEX_SHADER)   // Create a vertex shader object
	context.Call("shaderSource", vshader, vshader_source)      // Attach vertex shader source code
	context.Call("compileShader", vshader)                     // Compile the vertex shader
	fshader := context.Call("createShader", c.FRAGMENT_SHADER) // Create fragment shader object
	context.Call("shaderSource", fshader, fshader_source)      // Attach fragment shader source code
	context.Call("compileShader", fshader)                     // Compile the fragment shader
	shaderProgram := context.Call("createProgram")             // Create a shader program to combine the two shaders
	context.Call("attachShader", shaderProgram, vshader)       // Attach the compiled vertex shader
	context.Call("attachShader", shaderProgram, fshader)       // Attach the compiled fragment shader
	context.Call("linkProgram", shaderProgram)                 // Make the shader program linked
	context.Call("useProgram", shaderProgram)                  // Let the completed shader program to be used
	context.Call("deleteShader", vshader)
	context.Call("deleteShader", fshader)

	//// Attributes ////
	loc := context.Call("getAttribLocation", shaderProgram, "xyz")    // Get the location of attribute 'xyz' in the shader
	context.Call("vertexAttribPointer", loc, 3, c.FLOAT, false, 0, 0) // Point 'xyz' location to the positions of ARRAY_BUFFER
	context.Call("enableVertexAttribArray", loc)                      // Enable the use of attribute 'xyz' from ARRAY_BUFFER

	//// Draw the scene ////
	context.Call("clearColor", 1.0, 1.0, 1.0, 1.0) // Set clearing color
	context.Call("clear", c.COLOR_BUFFER_BIT)      // Clear the canvas
	context.Call("enable", c.DEPTH_TEST)           // Enable the depth test

	//// Draw the geometry ////
	context.Call("drawElements", c.TRIANGLES, len(indices), c.UNSIGNED_SHORT, 0)
	wcanvas.Run(nil)
}
