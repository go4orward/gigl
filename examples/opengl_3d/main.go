package main

import (
	"errors"
	"log"
	"runtime"

	"github.com/go4orward/gigl/env/opengl41"
	"github.com/go4orward/gigl/g3d"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	window, err := opengl41.NewOpenGLWindow(800, 600, "OpenGL3D: Rectangle with Gopher Texture Image", false)
	if err != nil {
		log.Fatal(errors.New("Failed to create OpenGL window : " + err.Error()))
	}
	rc := window.GetGLRenderingContext()
	scene := g3d.NewScene("#ffffff") // Scene with WHITE background

	geometry := g3d.NewGeometry_CubeWithTexture(1, 1, 1)    // create geometry (a cube of size 1.0)
	geometry.BuildNormalsForFace()                          // calculate normal vectors for each face
	geometry.BuildDataBuffers(true, false, true)            // build data buffers for vertices and faces
	material, _ := rc.CreateMaterial("./assets/gopher.png") // create material (with texture image)
	shader := g3d.NewShader_NormalTexture(rc)               // use the standard NORMAL+TEXTURE shader

	scene.Add(g3d.NewSceneObject(geometry, material, nil, nil, shader))
	// llayer := g3d.NewOverlayLabelLayer(rc, 20, true).AddLabelsForTest()
	// mlayer := g3d.NewOverlayMarkerLayer(rc).AddMarkersForTest()
	// scene.AddOverlay(mlayer)
	cam_ip := g3d.CamInternalParams{WH: rc.GetWH(), Fov: 15, Zoom: 1.0, NearFar: [2]float32{1, 100}}
	cam_ep := g3d.CamExternalPose{From: [3]float32{0, 0, 10}, At: [3]float32{0, 0, 0}, Up: [3]float32{0, 1, 0}}
	camera := g3d.NewCamera(true, &cam_ip, &cam_ep)
	renderer := g3d.NewRenderer(rc) // set up the renderer

	window.Run(func(now float64) {
		renderer.Clear(scene)               // prepare to render (clearing to white background)
		renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
		renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)
		scene.Get(0).Rotate([3]float32{0, 1, 1}, 1.0)
	})
}
