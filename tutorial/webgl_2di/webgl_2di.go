package main

import (
	"fmt"

	"github.com/go4orward/gigl/env/webgl10"
	"github.com/go4orward/gigl/g2d"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/webgl2dui_example.go'.
	fmt.Println("Hello WebGL 1.0")                      // printed in the browser console
	canvas, err := webgl10.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		fmt.Printf("Failed to start WebGL : %v\n", err)
		return
	}
	rc := canvas.GetRenderingContext()
	scene := g2d.NewScene("#ffffff")                            // Scene with WHITE background
	scene.Add(g2d.NewSceneObject_RectangleInstancesExample(rc)) // multiple instances of rectangles
	mlayer := g2d.NewOverlayMarkerLayer(rc).AddMarkersForTest()
	llayer := g2d.NewOverlayLabelLayer(rc, 20, true).AddLabelsForTest()
	scene.AddOverlay(mlayer, llayer)
	bbox, size, center := scene.GetBBoxSizeCenter(true)   // BoundingBox, Size(W&H) of BBox, Center of BBox
	camera := g2d.NewCamera(rc.GetWH(), size[0]*1.1, 1.0) // FOV covers the Width of BBox, ZoomLevel is 1.0
	camera.SetPose(center[0], center[1], 0.0).SetBoundingBox(bbox)
	renderer := g2d.NewRenderer(rc) // set up the renderer

	// set up user interactions
	canvas.SetEventHandlerForClick(func(canvasxy [2]int, keystat [4]bool) {
		wxy := camera.UnprojectCanvasToWorld(canvasxy)
		fmt.Printf("canvas (%d %d)  world (%.2f %.2f)\n", canvasxy[0], canvasxy[1], wxy[0], wxy[1])
	})
	canvas.SetEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
		fmt.Println(camera.Summary())
	})
	canvas.SetEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
		wdxy := camera.UnprojectCanvasDeltaToWorld(dxy)
		camera.Translate(-wdxy[0], -wdxy[1]).ApplyBoundingBox(true, false)
	})
	canvas.SetEventHandlerForZoom(func(canvasxy [2]int, scale float32, keystat [4]bool) {
		oldxy := camera.UnprojectCanvasToWorld(canvasxy)
		camera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
		newxy := camera.UnprojectCanvasToWorld(canvasxy)
		delta := g2d.NewV2dBySub(newxy, oldxy)
		camera.Translate(-delta[0], -delta[1]).ApplyBoundingBox(true, true)
	})
	canvas.SetEventHandlerForScroll(func(canvasxy [2]int, dx int, dy int, keystat [4]bool) {
		wdxy := camera.UnprojectCanvasDeltaToWorld([2]int{dx, dy})
		camera.Translate(0.0, wdxy[1]).ApplyBoundingBox(true, false)
	})
	canvas.SetEventHandlerForWindowResize(func(w int, h int) {
		camera.SetAspectRatio(w, h)
	})
	fmt.Println("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console

	// Run UI animation loop
	canvas.Run(func(now float64) {
		renderer.Clear(scene)               // prepare to render (clearing to white background)
		renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
		renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)
	})
}
