package g2d

import (
	"fmt"

	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g2d/c2d"
)

type SceneObject struct {
	Geometry    *Geometry         // geometry interface
	Material    common.GLMaterial // material
	VShader     common.GLShader   // vert shader and its bindings
	EShader     common.GLShader   // edge shader and its bindings
	FShader     common.GLShader   // face shader and its bindings
	modelmatrix common.Matrix3    // model transformation matrix of this SceneObject
	UseDepth    bool              // depth test flag (default is true)
	UseBlend    bool              // blending flag with alpha (default is false)
	poses       *SceneObjectPoses // OPTIONAL, poses for multiple instances of this (geometry+material) object
	children    []*SceneObject    // OPTIONAL, children of this SceneObject (to be rendered recursively)
	bbox        [2][2]float32     // bounding box
}

func NewSceneObject(geometry *Geometry, material common.GLMaterial,
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
	sobj := SceneObject{Geometry: geometry, Material: material, VShader: vshader, EShader: eshader, FShader: fshader}
	sobj.modelmatrix.SetIdentity()
	sobj.UseDepth = false // new drawings will overwrite old ones by default
	sobj.UseBlend = false // alpha blending is turned off by default
	sobj.poses = nil      // OPTIONAL, only if multiple instances of the geometry are rendered
	sobj.children = nil   // OPTIONAL, only if current SceneObject has any child SceneObjects
	sobj.bbox = c2d.BBoxInit()
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

func (self *SceneObject) SetTransformation(txy [2]float32, angle_in_degree float32, sxy [2]float32) *SceneObject {
	translation := common.NewMatrix3().SetTranslation(txy[0], txy[1])
	rotation := common.NewMatrix3().SetRotation(angle_in_degree)
	scaling := common.NewMatrix3().SetScaling(sxy[0], sxy[1])
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *SceneObject) Rotate(angle_in_degree float32) *SceneObject {
	rotation := common.NewMatrix3().SetRotation(angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Translate(tx float32, ty float32) *SceneObject {
	translation := common.NewMatrix3().SetTranslation(tx, ty)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *SceneObject) Scale(sx float32, sy float32) *SceneObject {
	scaling := common.NewMatrix3().SetScaling(sx, sy)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
	return self
}

// ----------------------------------------------------------------------------
// Bounding Box
// ----------------------------------------------------------------------------

func (self *SceneObject) GetBoundingBox(m *common.Matrix3, renew bool) [2][2]float32 {
	if !c2d.BBoxIsSet(self.bbox) || renew {
		bbox := c2d.BBoxInit()
		// apply the transformation matrx
		var mm *common.Matrix3 = nil
		if m != nil {
			mm = m.MultiplyToTheRight(&self.modelmatrix)
		} else {
			mm = self.modelmatrix.Copy() // new matrix
		}
		// add all the vertices of the geometry
		if self.poses == nil {
			for _, v := range self.Geometry.verts {
				xy := mm.MultiplyVector2(v)
				c2d.BBoxAddPoint(&bbox, xy)
			}
		} else {
			for i := 0; i < self.poses.Count; i++ {
				idx := i * self.poses.Size
				txy := self.poses.DataBuffer[idx : idx+2]
				for _, v := range self.Geometry.verts {
					xy := [2]float32{v[0] + txy[0], v[1] + txy[1]}
					xy = mm.MultiplyVector2(xy)
					c2d.BBoxAddPoint(&bbox, xy)
				}
			}
		}
		for _, sobj := range self.children {
			bbox = c2d.BBoxMerge(bbox, sobj.GetBoundingBox(mm, renew))
		}
		self.bbox = bbox
	}
	return self.bbox
}
