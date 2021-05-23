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
	c := rc.GetConstants()
	wrc := wcanvas.GetWebGLRenderingContext() // just for WebGL demo only

	// Build Geometry and its data buffers
	vertices := []float32{-0.5, 0.5, 0, -0.5, -0.5, 0, 0.5, -0.5, 0}
	indices := []uint32{2, 1, 0}
	var vertices_array = wcanvas.ConvertGoSliceToJsTypedArray(vertices)
	var indices_array = wcanvas.ConvertGoSliceToJsTypedArray(indices)
	vertexBuffer := wrc.Call("createBuffer", c.ARRAY_BUFFER)                     // create buffer
	wrc.Call("bindBuffer", c.ARRAY_BUFFER, vertexBuffer)                         // bind the buffer
	wrc.Call("bufferData", c.ARRAY_BUFFER, vertices_array, c.STATIC_DRAW)        // pass data to buffer
	indexBuffer := wrc.Call("createBuffer", c.ELEMENT_ARRAY_BUFFER)              // create index buffer
	wrc.Call("bindBuffer", c.ELEMENT_ARRAY_BUFFER, indexBuffer)                  // bind the buffer
	wrc.Call("bufferData", c.ELEMENT_ARRAY_BUFFER, indices_array, c.STATIC_DRAW) // pass data to the buffer

	// Shaders
	vshader_source := `
		attribute vec3 xyz;
		void main(void) {
			gl_Position = vec4(xyz, 1.0);
		}`
	fshader_source := `
		void main(void) {
			gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
		}`
	vshader := wrc.Call("createShader", c.VERTEX_SHADER)   // Create a vertex shader object
	wrc.Call("shaderSource", vshader, vshader_source)      // Attach vertex shader source code
	wrc.Call("compileShader", vshader)                     // Compile the vertex shader
	fshader := wrc.Call("createShader", c.FRAGMENT_SHADER) // Create fragment shader object
	wrc.Call("shaderSource", fshader, fshader_source)      // Attach fragment shader source code
	wrc.Call("compileShader", fshader)                     // Compile the fragment shader
	shaderProgram := wrc.Call("createProgram")             // Create a shader program to combine the two shaders
	wrc.Call("attachShader", shaderProgram, vshader)       // Attach the compiled vertex shader
	wrc.Call("attachShader", shaderProgram, fshader)       // Attach the compiled fragment shader
	wrc.Call("linkProgram", shaderProgram)                 // Make the shader program linked
	wrc.Call("useProgram", shaderProgram)                  // Let the completed shader program to be used
	wrc.Call("deleteShader", vshader)
	wrc.Call("deleteShader", fshader)

	// Bind Attributes with the data buffers
	loc := wrc.Call("getAttribLocation", shaderProgram, "xyz")    // Get the location of attribute 'xyz' in the shader
	wrc.Call("vertexAttribPointer", loc, 3, c.FLOAT, false, 0, 0) // Point 'xyz' location to the positions of ARRAY_BUFFER
	wrc.Call("enableVertexAttribArray", loc)                      // Enable the use of attribute 'xyz' from ARRAY_BUFFER

	// Prepare to draw
	wrc.Call("clearColor", 1.0, 1.0, 1.0, 1.0) // Set clearing color
	wrc.Call("clear", c.COLOR_BUFFER_BIT)      // Clear the canvas
	wrc.Call("enable", c.DEPTH_TEST)           // Enable the depth test

	// Draw the geometry
	wrc.Call("drawElements", c.TRIANGLES, len(indices), c.UNSIGNED_SHORT, 0)
	wcanvas.Run(nil)
}
