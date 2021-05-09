// +build js,wasm
package main

import (
	"fmt"

	"github.com/go4orward/gigl/common/webgl"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/webgl_example.go'.
	fmt.Println("Hello WebGL!")                        // printed in the browser console
	wcanvas, err := webgl.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		fmt.Printf("Failed to start WebGL : %v\n", err)
		return
	}
	vertices := []float32{-0.5, 0.5, 0, -0.5, -0.5, 0, 0.5, -0.5, 0}
	indices := []uint32{2, 1, 0}
	vertex_shader_code := `
		attribute vec3 xyz;
		void main(void) {
			gl_Position = vec4(xyz, 1.0);
		}`
	fragment_shader_code := `
		void main(void) {
			gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
		}`
	rc := wcanvas.GetRenderingContext()
	context, c := rc.GetContext(), rc.GetConstants()

	//// Geometry ////
	var vertices_array = rc.ConvertGoSliceToJsTypedArray(vertices)
	var indices_array = rc.ConvertGoSliceToJsTypedArray(indices)
	vertexBuffer := context.Call("createBuffer", c.ARRAY_BUFFER)                     // create buffer
	context.Call("bindBuffer", c.ARRAY_BUFFER, vertexBuffer)                         // bind the buffer
	context.Call("bufferData", c.ARRAY_BUFFER, vertices_array, c.STATIC_DRAW)        // pass data to buffer
	indexBuffer := context.Call("createBuffer", c.ELEMENT_ARRAY_BUFFER)              // create index buffer
	context.Call("bindBuffer", c.ELEMENT_ARRAY_BUFFER, indexBuffer)                  // bind the buffer
	context.Call("bufferData", c.ELEMENT_ARRAY_BUFFER, indices_array, c.STATIC_DRAW) // pass data to the buffer

	//// Shaders ////
	vertShader := context.Call("createShader", c.VERTEX_SHADER)    // Create a vertex shader object
	context.Call("shaderSource", vertShader, vertex_shader_code)   // Attach vertex shader source code
	context.Call("compileShader", vertShader)                      // Compile the vertex shader
	fragShader := context.Call("createShader", c.FRAGMENT_SHADER)  // Create fragment shader object
	context.Call("shaderSource", fragShader, fragment_shader_code) // Attach fragment shader source code
	context.Call("compileShader", fragShader)                      // Compile the fragment shader
	shaderProgram := context.Call("createProgram")                 // Create a shader program to combine the two shaders
	context.Call("attachShader", shaderProgram, vertShader)        // Attach the compiled vertex shader
	context.Call("attachShader", shaderProgram, fragShader)        // Attach the compiled fragment shader
	context.Call("linkProgram", shaderProgram)                     // Make the shader program linked
	context.Call("useProgram", shaderProgram)                      // Let the completed shader program to be used

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
}
