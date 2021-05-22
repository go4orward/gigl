package g3d

import (
	"math"

	"github.com/go4orward/gigl"
)

func NewSceneObject_CylinderWireframe(rc gigl.GLRenderingContext) *SceneObject {
	// This example creates a cylinder, to be rendered as 'wireframe'
	// (This example demonstrates how 'triangulation of face' works)
	geometry := NewGeometry_Cylinder(6, 0.5, 1.0, 0, true)      // create a cylinder with radius 0.5 and heigt 1.0
	geometry.BuildDataBuffersForWireframe()                     // extract wireframe edges from faces
	material, _ := rc.CreateMaterial("#888888")                 // create material with yellow color
	shader := NewShader_ColorOnly(rc)                           // use the standard COLOR_ONLY shader
	return NewSceneObject(geometry, material, nil, shader, nil) // set up the scene object (draw EDGES only)
}

func NewSceneObject_CubeWithTexture(rc gigl.GLRenderingContext) *SceneObject {
	geometry := NewGeometry_CubeWithTexture(1.0, 1.0, 1.0)
	geometry.BuildNormalsForFace()
	geometry.BuildDataBuffers(true, false, true)                // build data buffers for vertices and faces
	material, _ := rc.CreateMaterial("/assets/gopher.png")      // create material with a texture image
	shader := NewShader_NormalTexture(rc)                       // use the standard NORMAL+TEXTURE shader
	return NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object (draw FACES only)
}

func NewSceneObject_CubeInstances(rc gigl.GLRenderingContext) *SceneObject {
	// This example creates 40,000 instances of a single geometry, each with its own pose (tx, ty)
	geometry := NewGeometry_Cube(0.08, 0.08, 0.08)                 // create a cube of size 0.08
	geometry.BuildNormalsForFace()                                 // prepare face normal vectors
	geometry.BuildDataBuffers(true, false, true)                   //
	material, _ := rc.CreateMaterial("#888888")                    // create material
	shader := NewShader_InstancePoseColor(rc)                      // create shader, and set its bindings
	scnobj := NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object (draw FACES only)
	scnobj.SetInstanceBuffer(10*10*10, 6, nil)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 10; k++ {
				scnobj.SetInstancePoseValues(i*100+j*10+k, 0, float32(i)/10, float32(j)/10, float32(k)/10) // tx, ty, tz
				ii, jj, kk := math.Abs(float64(i)-5)/5, math.Abs(float64(j)-5)/5, math.Abs(float64(k)-5)/5
				r, g, b := float32(ii), float32(jj), float32(kk)
				// r, g, b := uint8(float32(ii)*255), uint8(float32(jj)*255), uint8(float32(kk)*255)
				scnobj.SetInstancePoseValues(i*100+j*10+k, 3, r, g, b, 255) // color
			}
		}
	}
	scnobj.Translate(-0.5, -0.5, -0.5)
	return scnobj
}

func NewSceneObject_Airplane(rc gigl.GLRenderingContext) *SceneObject {
	centers := [][3]float32{{0, -0.025, -1}, {0, -0.025, -0.99}, {0, -0.02, -0.9}, {0, -0.01, -0.6}, {0, 0, +0.0}, {0, 0, +0.8}, {0, 0, +0.9}, {0, 0, +0.99}, {0, 0, +1}}
	radii := []float32{0, 0.01, 0.04, 0.08, 0.1, 0.1, 0.08, 0.02, 0}
	wingth := float32(0.02)
	pbody := NewGeometry_SolidFromCentersAndRadii(centers, radii, 8)
	rwing := NewGeometry_SolidFromFaceAndHeight([][3]float32{{0, 0.0}, {+0.8, -0.3}, {+0.8, -0.2}, {0, 0.4}}, wingth)
	lwing := NewGeometry_SolidFromFaceAndHeight([][3]float32{{0, 0.4}, {-0.8, -0.2}, {-0.8, -0.3}, {0.0, 0}}, wingth)
	twing := NewGeometry_SolidFromFaceAndHeight([][3]float32{{0, 0.3}, {-0.3, 0.05}, {-0.3, 0}, {+0.3, 0}, {+0.3, 0.05}}, wingth)
	vwing := NewGeometry_SolidFromFaceAndHeight([][3]float32{{0, 0.0}, {0.3, 0}, {0.3, 0.05}, {0, 0.3}}, wingth).Rotate([3]float32{0, 1, 0}, -90)
	geometry := NewGeometry()
	geometry.Merge(pbody.Rotate([3]float32{1, 0, 0}, -90))
	geometry.Merge(rwing.Translate(0, 0.0, -wingth))
	geometry.Merge(lwing.Translate(0, 0.0, -wingth))
	geometry.Merge(twing.Translate(0, -0.9, +wingth))
	geometry.Merge(vwing.Translate(wingth/2, -0.9, 0))
	geometry.BuildNormalsForVertex()                               // prepare normal vectors
	geometry.BuildDataBuffers(true, false, true)                   //
	material, _ := rc.CreateMaterial("#ffff88")                    // create material
	shader := NewShader_NormalColor(rc)                            // use the standard NORMAL+COLOR shader
	scnobj := NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object (draw FACES only)
	return scnobj
}
