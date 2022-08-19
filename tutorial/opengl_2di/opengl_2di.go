package main

import (
	"fmt"

	"github.com/go4orward/gigl/common"
	opengl "github.com/go4orward/gigl/env/opengl41"
	"github.com/go4orward/gigl/g2d"
)

type Config struct {
	loglevel  string //
	logfilter string //
}

func main() {
	cfg := Config{loglevel: "trace", logfilter: ""}
	if cfg.loglevel != "" {
		common.SetLogger(common.NewConsoleLogger(cfg.loglevel)).SetTraceFilter(cfg.logfilter).SetOption("", false)
	}
	canvas, err := opengl.NewOpenGLCanvas(1200, 900, "OpenGL2D: Triangle in World Space", false)
	if err != nil {
		common.Logger.Error("Failed to start WebGL : %v\n", err)
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

	SetUIEventHandlers(canvas, camera)
	common.Logger.Trace("SceneObject\n%s", scene.Get(0).Summary())

	// Run UI animation loop
	canvas.Run(func(now float64) {
		renderer.Clear(scene)               // prepare to render (clearing to white background)
		renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
		renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)
	})
}

func SetUIEventHandlers(canvas *opengl.OpenGLCanvas, camera *g2d.Camera) {
	// set up user interactions
	canvas.SetEventHandlerForClick(func(canvasxy [2]int, keystat [4]bool) {
		wxy := camera.UnprojectCanvasToWorld(canvasxy)
		common.Logger.Info("canvas (%d %d)  world (%.2f %.2f)\n", canvasxy[0], canvasxy[1], wxy[0], wxy[1])
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
	common.Logger.Info("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console
}
