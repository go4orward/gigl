package main

import (
	"fmt"

	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/env/webgl10"
	webgl "github.com/go4orward/gigl/env/webgl10"
	"github.com/go4orward/gigl/g2d"
	"github.com/go4orward/gigl/g3d"
)

type Config struct {
	loglevel  string //
	logfilter string //
}

func main() {
	cfg := Config{loglevel: "info", logfilter: ""}
	if cfg.loglevel != "" {
		common.SetLogger(common.NewConsoleLogger(cfg.loglevel)).SetTraceFilter(cfg.logfilter).SetOption("", false)
	}
	// THIS CODE IS SUPPOSED TO BE BUILT AS WEBASSEMBLY AND RUN INSIDE A BROWSER.
	// BUILD IT LIKE 'GOOS=js GOARCH=wasm go build -o example.wasm examples/example.go'.
	canvas, err := webgl.NewWebGLCanvas("wasmcanvas") // ID of canvas element
	if err != nil {
		common.Logger.Error("Failed to start WebGL : %v\n", err)
		return
	}
	rc := canvas.GetRenderingContext()
	scene := g3d.NewScene("#ffffff")                         // Scene with WHITE background
	geometry := g3d.NewGeometryCubeWithTexture(1, 1, 1)      // create geometry (a cube of size 1.0)
	geometry.BuildNormalsForFace()                           // calculate normal vectors for each face
	geometry.BuildDataBuffers(true, false, true)             // build data buffers for vertices and faces
	material := g2d.NewMaterialTexture("/assets/gopher.png") // create material (with texture image)
	shader := g3d.NewShader_NormalTexture(rc)                // use the standard NORMAL+TEXTURE shader
	scene.Add(g3d.NewSceneObject(geometry, material, nil, nil, shader))
	// llayer := g3d.NewOverlayLabelLayer(rc, 20, true).AddLabelsForTest()
	// mlayer := g3d.NewOverlayMarkerLayer(rc).AddMarkersForTest()
	// scene.AddOverlay(mlayer, llayer)
	cam_ip := g3d.CamInternalParams{WH: rc.GetWH(), Fov: 15, Zoom: 1.0, NearFar: [2]float32{1, 100}}
	cam_ep := g3d.CamExternalPose{From: [3]float32{0, 0, 10}, At: [3]float32{0, 0, 0}, Up: [3]float32{0, 1, 0}}
	camera := g3d.NewCamera(true, &cam_ip, &cam_ep)
	renderer := g3d.NewRenderer(rc) // set up the renderer
	//
	if common.Logger.IsLogging(common.LogLevelTrace) {
		common.Logger.Trace("SceneObject \n%s", scene.Get(0).Summary())
	}

	SetUIEventHandlers(canvas, camera)

	// run UI animation loop
	canvas.Run(func(now float64) {
		renderer.Clear(scene)               // prepare to render (clearing to white background)
		renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
		renderer.RenderAxes(camera, 0.8)    // render the axes (just for visual reference)
		scene.Get(0).Rotate([3]float32{0, 1, 1}, 1.0)
	})
}

func SetUIEventHandlers(canvas *webgl10.WebGLCanvas, camera *g3d.Camera) {
	// set up user interactions
	canvas.SetEventHandlerForClick(func(canvasxy [2]int, keystat [4]bool) {
		common.Logger.Info("%v\n", canvasxy)
	})
	canvas.SetEventHandlerForDoubleClick(func(canvasxy [2]int, keystat [4]bool) {
		fmt.Println(camera.Summary())
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
		switch code {
		default:
			common.Logger.Info("keypress : key='%v' code='%v'  keystat=%v\n", key, code, keystat)
		}
	})
	common.Logger.Info("Try mouse drag & wheel with SHIFT key pressed") // printed in the browser console
}
