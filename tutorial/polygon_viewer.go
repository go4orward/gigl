package main

import (
	"errors"
	"log"

	"github.com/go4orward/gigl/env/opengl41"
)

func main() {
	window, err := opengl41.NewOpenGLWindow(800, 600, "OpenGL3D: Rectangle with Gopher Texture Image", false)
	if err != nil {
		log.Fatal(errors.New("Failed to create OpenGL window : " + err.Error()))
	}
	scene := g3d.scene.NewScene(window.GetRenderingContext())

	geometry := g3d.geometry.NewCubeWithTexture(1, 1, 1)          // create geometry (a cube of size 1.0)
	geometry.BuildNormalsForFace()                                // calculate normal vectors for each face
	geometry.BuildDataBuffers(true, false, true)                  // build data buffers for vertices and faces
	material, _ := g3d.material.NewTexture("./assets/gopher.png") // create material (with texture image)

	// shader := g3d.NewShader_NormalTexture(rc)               // use the standard NORMAL+TEXTURE shader
	scene.Add(g3d.scene.NewObject(geometry, material, nil, nil))
	// llayer := g3d.NewOverlayLabelLayer(rc, 20, true).AddLabelsForTest()
	// mlayer := g3d.NewOverlayMarkerLayer(rc).AddMarkersForTest()
	// scene.AddOverlay(mlayer)
	cam_ip := g3d.camera.CamInternalParams{WH: window.GetWH(), Fov: 15, Zoom: 1.0, NearFar: [2]float32{1, 100}}
	cam_ep := g3d.camera.CamExternalPose{From: [3]float32{0, 0, 10}, At: [3]float32{0, 0, 0}, Up: [3]float32{0, 1, 0}}
	camera := g3d.camera.NewPerspective(true, &cam_ip, &cam_ep)
	// renderer := g3d.NewRenderer(rc) // set up the renderer
	camera.SetKeyFrames()

	// window.Run(func(now float64) {
	// 	renderer.Clear(scene)               // prepare to render (clearing to white background)
	// 	renderer.RenderScene(scene, camera) // render the scene (iterating over all the SceneObjects in it)
	// 	renderer.RenderAxes(camera, 1.0)    // render the axes (just for visual reference)
	// 	scene.Get(0).Rotate([3]float32{0, 1, 1}, 1.0)
	// })
	window.Play(scene, camera)
}
