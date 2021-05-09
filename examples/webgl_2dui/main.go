// +build js,wasm
package main

import (
	"fmt"

	"github.com/go4orward/gigl/env/webgl"
	"github.com/go4orward/gigl/g2d"
	"github.com/go4orward/gigl/g2d/c2d"
)

func main() {
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/webgl2dui_example.go'.
	fmt.Println("Hello WebGL!")                        // printed in the browser console
	wcanvas, err := webgl.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		fmt.Printf("Failed to start WebGL : %v\n", err)
		return
	}
	rc := wcanvas.GetRenderingContext()
	scene := g2d.NewScene("#ffffff") // Scene with WHITE background
	if true {
		scene.Add(g2d.NewSceneObject_RectangleInstancesExample(rc)) // multiple instances of rectangles
		mlayer := g2d.NewOverlayMarkerLayer(rc).AddMarkersForTest()
		llayer := g2d.NewOverlayLabelLayer(rc, 20, true).AddLabelsForTest()
		scene.AddOverlay(mlayer, llayer)
	} else {
		geometry := g2d.NewGeometry_Rectangle(1.0) // create geometry (a rectangle)
		geometry.SetTextureUVs([][]float32{{0, 1}, {1, 1}, {1, 0}, {0, 0}})
		geometry.BuildDataBuffers(true, true, true)            // build data buffers for vertices, edges, and faces
		material, _ := rc.CreateMaterial("/assets/gopher.png") // create material (with texture image)
		shader := g2d.NewShader_MaterialTexture(rc)            // shader with auto-binded color & PVM matrix
		scnobj := g2d.NewSceneObject(geometry, material, nil, nil, shader)
		scene.Add(scnobj)
	}
	bbox, size, center := scene.GetBBoxSizeCenter(true)   // BoundingBox, Size(W&H) of BBox, Center of BBox
	camera := g2d.NewCamera(rc.GetWH(), size[0]*1.1, 1.0) // FOV covers the Width of BBox, ZoomLevel is 1.0
	camera.SetPose(center[0], center[1], 0.0).SetBoundingBox(bbox)
	renderer := g2d.NewRenderer(rc)     // set up the renderer
	renderer.Clear(scene)               // prepare to render (clearing to white background)
	renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
	renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)

	if true { // ONLY FOR INTERACTIVE UI
		fmt.Println("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console
		// add user interactions (with mouse)
		wcanvas.SetupEventHandlers()
		wcanvas.RegisterEventHandlerForClick(func(canvasxy [2]int, keystat [4]bool) {
			wxy := camera.UnprojectCanvasToWorld(canvasxy)
			fmt.Printf("canvas (%d %d)  world (%.2f %.2f)\n", canvasxy[0], canvasxy[1], wxy[0], wxy[1])
		})
		wcanvas.RegisterEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
			camera.ShowInfo()
		})
		wcanvas.RegisterEventHandlerForMouseDrag(func(canvasxy [2]int, dxy [2]int, keystat [4]bool) {
			wdxy := camera.UnprojectCanvasDeltaToWorld(dxy)
			camera.Translate(-wdxy[0], -wdxy[1]).ApplyBoundingBox(true, false)
		})
		wcanvas.RegisterEventHandlerForMouseWheel(func(canvasxy [2]int, scale float32, keystat [4]bool) {
			if keystat[3] { // ZOOM, if SHIFT key was pressed (ALT,CTRL,META,SHIFT)
				oldxy := camera.UnprojectCanvasToWorld(canvasxy)
				camera.SetZoom(scale) // 'scale' in [ 0.01 ~ 1(default) ~ 100.0 ]
				newxy := camera.UnprojectCanvasToWorld(canvasxy)
				delta := c2d.SubAB(newxy, oldxy)
				camera.Translate(-delta[0], -delta[1]).ApplyBoundingBox(true, true)
			} else { // SCROLL
				wdxy := camera.UnprojectCanvasDeltaToWorld([2]int{0, int(scale)})
				camera.Translate(0.0, wdxy[1]).ApplyBoundingBox(true, false)
			}
		})
		wcanvas.RegisterEventHandlerForWindowResize(func(w int, h int) {
			camera.SetAspectRatio(w, h)
		})
		// add animation
		wcanvas.SetupAnimationFrame(func() {
			renderer.Clear(scene)               // prepare to render (clearing to white background)
			renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
			renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)
		})
	}
	<-make(chan bool) // wait for events (without exiting)
}
