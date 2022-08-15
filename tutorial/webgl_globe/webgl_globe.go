package main

import (
	"fmt"

	"github.com/go4orward/gigl/earth"
	"github.com/go4orward/gigl/env/webgl10"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/webglglobe_example.go'
	fmt.Println("Hello WebGL 1.0")                      // printed in the browser console
	canvas, err := webgl10.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		fmt.Printf("Failed to start WebGL : %v\n", err)
		return
	}
	rc := canvas.GetRenderingContext()
	wglobe := earth.NewWorldGlobe(rc, "#000000", "/assets/world.png") // Globe radius is assumed to be 1.0
	wcamera := earth.NewWorldCamera(rc.GetWH(), 15, 1.0)              // camera FOV default is 15° (in degree)
	wcamera.SetPoseByLonLat(0, 0, 10)                                 // longitude 0°, latitude 0°, radius(distance) 10.0
	renderer := earth.NewWorldRenderer(rc)                            // set up the world renderer

	// add user interactions (with mouse)
	canvas.SetEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
		fmt.Printf("%s\n", wcamera.Summary())
	})
	canvas.SetEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
		wcamera.RotateAroundGlobe(float32(dxy[0])*0.2, float32(dxy[1])*0.2)
	})
	canvas.SetEventHandlerForZoom(func(canvasxy [2]int, scale float32, keystat [4]bool) {
		wcamera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
	})
	canvas.SetEventHandlerForWindowResize(func(w int, h int) {
		wcamera.SetAspectRatio(w, h)
	})
	fmt.Println("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console

	// run UI animation loop
	canvas.Run(func(now float64) {
		renderer.Clear(wglobe)                  // prepare to render (clearing to black background)
		renderer.RenderWorld(wglobe, wcamera)   // render the Globe (and all the layers & glowring)
		wglobe.Rotate([3]float32{0, 0, 1}, 0.1) // rotate the Globe
	})
}
