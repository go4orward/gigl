package g3d

import (
	"fmt"

	"github.com/go4orward/gigl/common"
)

type SceneObject struct {
	Geometry    common.GLGeometry // geometry interface
	Material    common.GLMaterial // material
	VShader     common.GLShader   // vert shader and its bindings
	EShader     common.GLShader   // edge shader and its bindings
	FShader     common.GLShader   // face shader and its bindings
	modelmatrix common.Matrix4    //
	UseDepth    bool              // depth test flag (default is true)
	UseBlend    bool              // blending flag with alpha (default is false)
	poses       *SceneObjectPoses // poses for multiple instances of this (geometry+material) object
	children    []*SceneObject    //
}

func NewSceneObject(geometry common.GLGeometry, material common.GLMaterial,
	vshader common.GLShader, eshader common.GLShader, fshader common.GLShader) *SceneObject {
	// 'geometry' : geometric shape (vertices, edges, faces) to be rendered
	// 'material' : color, texture, or other material properties	: OPTIONAL (can be 'nil')
	// 'vshader' : shader for VERTICES (POINTS) 					: OPTIONAL (can be 'nil')
	// 'eshader' : shader for EDGES (LINES) 						: OPTIONAL (can be 'nil')
	// 'fshader' : shader for FACES (TRIANGLES) 					: OPTIONAL (can be 'nil')
	// Note that geometry & material & shader can be shared among different SceneObjects.
	if geometry == nil {
		return nil
	}
	// Note that 'material' & 'shader' can be nil, in which case its parent's 'material' & 'shader' will be used to render.
	sobj := SceneObject{Geometry: geometry, Material: material, VShader: vshader, EShader: eshader, FShader: fshader}
	sobj.modelmatrix.SetIdentity()
	sobj.UseDepth = true  // depth test is turned on by default
	sobj.UseBlend = false // alpha blending is turned off by default
	sobj.poses = nil
	sobj.children = nil
	return &sobj
}

func (self *SceneObject) ShowInfo() {
	fmt.Printf("SceneObject ")
	self.Geometry.ShowInfo()
	if self.poses != nil {
		fmt.Printf("  ")
		self.poses.ShowInfo()
	}
	if self.Material != nil {
		fmt.Printf("  ")
		self.Material.ShowInfo()
	}
	if self.VShader != nil {
		fmt.Printf("  VERT ")
		self.VShader.ShowInfo()
	}
	if self.EShader != nil {
		fmt.Printf("  EDGE ")
		self.EShader.ShowInfo()
	}
	if self.FShader != nil {
		fmt.Printf("  FACE ")
		self.FShader.ShowInfo()
	}
	fmt.Printf("  Flags    : UseDepth=%t  UseBlend=%t\n", self.UseDepth, self.UseBlend)
	fmt.Printf("  Children : %d\n", len(self.children))
}

// ----------------------------------------------------------------------------
// Basic Access
// ----------------------------------------------------------------------------

func (self *SceneObject) AddChild(child *SceneObject) *SceneObject {
	if self.children == nil {
		self.children = make([]*SceneObject, 0)
	}
	self.children = append(self.children, child)
	return self
}

func (self *SceneObject) GetModelMatrix() *common.Matrix4 {
	return &self.modelmatrix
}

func (self *SceneObject) GetChildren() []*SceneObject {
	return self.children
}

// ----------------------------------------------------------------------------
// Multiple Instance Poses
// ----------------------------------------------------------------------------

func (self *SceneObject) SetPoses(poses *SceneObjectPoses) *SceneObject {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	self.poses = poses
	return self
}

func (self *SceneObject) SetupPoses(size int, count int, data []float32) *SceneObject {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	self.poses = NewSceneObjectPoses(size, count, data)
	return self
}

func (self *SceneObject) SetPoseValues(index int, offset int, values ...float32) *SceneObject {
	// This function is OPTIONAL (only if multiple instances of the geometry are rendered)
	self.poses.SetPose(index, offset, values...)
	return self
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *SceneObject) SetTransformation(txyz [3]float32, axis [3]float32, angle_in_degree float32, sxyz [3]float32) *SceneObject {
	translation := common.NewMatrix4().SetTranslation(txyz[0], txyz[1], txyz[2])
	rotation := common.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	scaling := common.NewMatrix4().SetScaling(sxyz[0], sxyz[1], sxyz[2])
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *SceneObject) Translate(tx float32, ty float32, tz float32) *SceneObject {
	translation := common.NewMatrix4().SetTranslation(tx, ty, tz)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Rotate(axis [3]float32, angle_in_degree float32) *SceneObject {
	rotation := common.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Scale(sx float32, sy float32, sz float32) *SceneObject {
	scaling := common.NewMatrix4().SetScaling(sx, sy, sz)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
	return self
}
