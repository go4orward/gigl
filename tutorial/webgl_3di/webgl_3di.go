package main

import (
	"fmt"
	"math"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/env/webgl10"
	"github.com/go4orward/gigl/g2d"
	"github.com/go4orward/gigl/g3d"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/webgl3d_example.go'.
	fmt.Println("Hello WebGL 1.0")                      // printed in the browser console
	canvas, err := webgl10.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		fmt.Printf("Failed to start WebGL : %v\n", err)
		return
	}
	rc := canvas.GetRenderingContext()

	// This example creates 40,000 instances of a single geometry, each with its own pose (tx, ty)
	geometry := g3d.NewGeometryCube(0.8, 0.8, 0.8)                     // create a cube of size 0.08
	geometry.BuildNormalsForFace()                                     // prepare face normal vectors
	geometry.BuildDataBuffers(true, false, true)                       //
	material := g2d.NewMaterialColors("#888888")                       // create material
	shader := get_shader_with_instance_pose_and_color(rc)              // create shader, and set its bindings
	scnobj := g3d.NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object (draw FACES only)
	N, hN := 50, 25
	scnobj.SetInstanceBuffer(N*N*N, 4, nil) // (3 for xyz + 1 for color)
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			for k := 0; k < N; k++ {
				x, y, z := float32(i-hN), float32(j-hN), float32(k-hN)
				scnobj.SetInstancePoseValues(i*N*N+j*N+k, 0, x, y, z) // cube position
				ii, jj, kk := float64(i-hN)/float64(hN), float64(j-hN)/float64(hN), float64(k-hN)/float64(hN)
				r, g, b := uint8(math.Abs(ii)*255), uint8(math.Abs(jj)*255), uint8(math.Abs(kk)*255)
				scnobj.SetInstanceColorValues(i*N*N+j*N+k, 3, r, g, b, 255) // cube color (put 4 * uint8 into 1 * float32)
			}
		}
	}
	scnobj.Translate(-0.5, -0.5, -0.5)

	scene := g3d.NewScene("#ffffff").Add(scnobj) // Scene with WHITE background
	cam_ip := g3d.CamInternalParams{WH: rc.GetWH(), Fov: 30, Zoom: 1.0, NearFar: [2]float32{1, 600}}
	cam_ep := g3d.CamExternalPose{From: [3]float32{0, 0, 100}, At: [3]float32{0, 0, 0}, Up: [3]float32{0, 1, 0}}
	camera := g3d.NewCamera(true, &cam_ip, &cam_ep)
	renderer := g3d.NewRenderer(rc) // set up the renderer

	// set up user interactions
	canvas.SetEventHandlerForClick(func(canvasxy [2]int, keystat [4]bool) {
		fmt.Printf("%v\n", canvasxy)
	})
	canvas.SetEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
		fmt.Printf("%v\n", camera.Summary())
	})
	canvas.SetEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
		camera.RotateAroundPoint(10, float32(dxy[0])*0.2, float32(dxy[1])*0.2)
	})
	canvas.SetEventHandlerForZoom(func(canvasxy [2]int, scale float32, keystat [4]bool) {
		camera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
	})
	canvas.SetEventHandlerForWindowResize(func(w int, h int) {
		camera.SetAspectRatio(w, h)
	})
	canvas.SetEventHandlerForKeyPress(func(key string, code string, keystat [4]bool) {
		if code == "Space" {
			fmt.Printf("keypress : %v\n", code)
		}
	})
	fmt.Println("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console

	// run UI animation loop
	canvas.Run(func(now float64) {
		renderer.Clear(scene)               // prepare to render (clearing to white background)
		renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
		renderer.RenderAxes(camera, 0.8)    // render the axes (just for visual reference)
		scene.Get(0).Rotate([3]float32{0, 1, 1}, 1.0)
	})
}

func get_shader_with_instance_pose_and_color(rc gigl.GLRenderingContext) gigl.GLShader {
	// Shader for (XYZ + NORMAL) Geometry & (COLOR) Material & (DIRECTIONAL) Lighting
	var vertex_shader_code = `
		precision mediump float;
		uniform mat4 proj;			// Projection matrix
		uniform mat4 vwmd;			// ModelView matrix
		attribute vec3 xyz;			// XYZ coordinates
		attribute vec3 nor;			// normal vector
		attribute vec3 ixyz;		// instance pose : XYZ translation
		attribute vec3 icolor;		// instance pose : color
		uniform mat3 light;			// [0]: direction, [1]: color, [2]: ambient_color   (column-major)
		varying vec3 v_color;    	// (varying) instance color
		varying vec3 v_light;    	// (varying) lighting intensity
		void main() {
			gl_Position = proj * vwmd * vec4(xyz.x + ixyz[0], xyz.y + ixyz[1], xyz.z + ixyz[2], 1.0);
			float s = sqrt( vwmd[0][0]*vwmd[0][0] + vwmd[0][1]*vwmd[0][1] + vwmd[0][2]*vwmd[0][2]);  // scaling 
			mat3  mvRot = mat3( vwmd[0][0]/s, vwmd[0][1]/s, vwmd[0][2]/s, vwmd[1][0]/s, vwmd[1][1]/s, vwmd[1][2]/s, vwmd[2][0]/s, vwmd[2][1]/s, vwmd[2][2]/s );
			vec3  normal    = mvRot * nor;               		// normal vector in camera space
			float intensity = max(dot(normal, light[0]), 0.0);	// light_intensity = dot(face_normal,light_direction)
			v_light = intensity * light[1] + light[2];        	// intensity * light_color + ambient_color
			v_color = icolor;
		}`
	var fragment_shader_code = `
		precision mediump float;
		varying vec3 v_color;		// (varying) instance color
		varying vec3 v_light;		// (varying) lighting intensity
		void main() { 
			gl_FragColor = vec4(v_color * v_light, 1.0);
		}`
	shader, _ := rc.CreateShader(vertex_shader_code, fragment_shader_code)
	shader.SetBindingForUniform("proj", "mat4", "renderer.proj")          // (Projection) matrix
	shader.SetBindingForUniform("vwmd", "mat4", "renderer.vwmd")          // (View * Models) matrix
	shader.SetBindingForUniform("light", "mat3", "lighting.dlight")       // directional lighting
	shader.SetBindingForAttribute("xyz", "vec3", "geometry.coords")       // point XYZ coordinates
	shader.SetBindingForAttribute("nor", "vec3", "geometry.normal")       // point normal vectors
	shader.SetBindingForAttribute("ixyz", "vec3", "instance.pose:4:0")    // instance position of xyz coordinates
	shader.SetBindingForAttribute("icolor", "vec3", "instance.color:4:3") // instance color (packed in 1 float32)
	shader.CheckBindings()                                                // check validity of the shader
	return shader
}
