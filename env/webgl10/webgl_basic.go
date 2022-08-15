package webgl10

func DrawSimplestTriangle(canvas *WebGLCanvas) {
	wrc := canvas.GetWebGLRenderingContext() // just for WebGL demo only
	wc := canvas.GetWebGLConstants()

	// Build Geometry and its data buffers
	vertices := []float32{-0.5, 0.5, 0, -0.5, -0.5, 0, 0.5, -0.5, 0}
	indices := []uint32{2, 1, 0}
	var vertices_array = canvas.ConvertGoSliceToJsTypedArray(vertices)
	var indices_array = canvas.ConvertGoSliceToJsTypedArray(indices)
	vertexBuffer := wrc.Call("createBuffer", wc.ARRAY_BUFFER)                      // create buffer
	wrc.Call("bindBuffer", wc.ARRAY_BUFFER, vertexBuffer)                          // bind the buffer
	wrc.Call("bufferData", wc.ARRAY_BUFFER, vertices_array, wc.STATIC_DRAW)        // pass data to buffer
	indexBuffer := wrc.Call("createBuffer", wc.ELEMENT_ARRAY_BUFFER)               // create index buffer
	wrc.Call("bindBuffer", wc.ELEMENT_ARRAY_BUFFER, indexBuffer)                   // bind the buffer
	wrc.Call("bufferData", wc.ELEMENT_ARRAY_BUFFER, indices_array, wc.STATIC_DRAW) // pass data to the buffer

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
	vshader := wrc.Call("createShader", wc.VERTEX_SHADER)   // Create a vertex shader object
	wrc.Call("shaderSource", vshader, vshader_source)       // Attach vertex shader source code
	wrc.Call("compileShader", vshader)                      // Compile the vertex shader
	fshader := wrc.Call("createShader", wc.FRAGMENT_SHADER) // Create fragment shader object
	wrc.Call("shaderSource", fshader, fshader_source)       // Attach fragment shader source code
	wrc.Call("compileShader", fshader)                      // Compile the fragment shader
	shaderProgram := wrc.Call("createProgram")              // Create a shader program to combine the two shaders
	wrc.Call("attachShader", shaderProgram, vshader)        // Attach the compiled vertex shader
	wrc.Call("attachShader", shaderProgram, fshader)        // Attach the compiled fragment shader
	wrc.Call("linkProgram", shaderProgram)                  // Make the shader program linked
	wrc.Call("useProgram", shaderProgram)                   // Let the completed shader program to be used
	wrc.Call("deleteShader", vshader)
	wrc.Call("deleteShader", fshader)

	// Bind Attributes with the data buffers
	loc := wrc.Call("getAttribLocation", shaderProgram, "xyz")     // Get the location of attribute 'xyz' in the shader
	wrc.Call("vertexAttribPointer", loc, 3, wc.FLOAT, false, 0, 0) // Point 'xyz' location to the positions of ARRAY_BUFFER
	wrc.Call("enableVertexAttribArray", loc)                       // Enable the use of attribute 'xyz' from ARRAY_BUFFER

	// Prepare to draw
	wrc.Call("clearColor", 1.0, 1.0, 1.0, 1.0) // Set clearing color
	wrc.Call("clear", wc.COLOR_BUFFER_BIT)     // Clear the canvas
	wrc.Call("enable", wc.DEPTH_TEST)          // Enable the depth test

	// Draw the geometry
	wrc.Call("drawElements", wc.TRIANGLES, len(indices), wc.UNSIGNED_SHORT, 0)

}
