package globe

import (
	"math"

	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g3d"
	"github.com/go4orward/gigl/g3d/c3d"
)

type Globe struct {
	bkgcolor    [3]float32       // background color of the globe
	GSphere     *g3d.SceneObject // globe sphere with texture & vertex normals
	GlowRing    *g3d.SceneObject // glow ring around the globe
	modelmatrix common.Matrix4   // Model matrix of the globe & its layers
}

func NewGlobe(rc common.GLRenderingContext, bkg_color string) *Globe {
	self := Globe{}
	self.SetBkgColor(bkg_color)
	self.GSphere = NewSceneObject_Globe(rc)     // globe sphere with texture & vertex normals
	self.GlowRing = NewSceneObject_GlowRing(rc) // glow ring around the globe
	self.modelmatrix.SetIdentity()              // initialize as Identity matrix
	return &self
}

func (self *Globe) IsReadyToRender() bool {
	return self.GSphere.Material.IsTextureLoading() == false
}

// ----------------------------------------------------------------------------
// Background Color
// ----------------------------------------------------------------------------

func (self *Globe) SetBkgColor(color string) *Globe {
	rgba := common.ParseHexColor(color)
	self.bkgcolor = [3]float32{rgba[0], rgba[1], rgba[2]}
	return self
}

func (self *Globe) GetBkgColor() [3]float32 {
	return self.bkgcolor
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *Globe) SetTransformation(txyz [3]float32, axis [3]float32, angle_in_degree float32, scale float32) *Globe {
	translation := common.NewMatrix4().SetTranslation(txyz[0], txyz[1], txyz[2])
	rotation := common.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	scaling := common.NewMatrix4().SetScaling(scale, scale, scale)
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *Globe) Translate(tx float32, ty float32, tz float32) *Globe {
	translation := common.NewMatrix4().SetTranslation(tx, ty, tz)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *Globe) Rotate(axis [3]float32, angle_in_degree float32) *Globe {
	rotation := common.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *Globe) Scale(scale float32) *Globe {
	scaling := common.NewMatrix4().SetScaling(scale, scale, scale)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
	return self
}

// ----------------------------------------------------------------------------
// Globe
// ----------------------------------------------------------------------------

func NewSceneObject_GlobeWithoutLight(rc common.GLRenderingContext) *g3d.SceneObject {
	// Globe model with texture UV coordinates (without normal vectors and directional lighting)
	geometry := build_globe_geometry(1.0, 64, 32, false)            // create globe geometry with texture UVs only
	geometry.BuildDataBuffers(true, false, true)                    // build data buffers for vertices and faces
	material, _ := rc.CreateMaterial("/assets/world.png")           // create material with a texture of world image
	shader := g3d.NewShader_TextureOnly(rc)                         // use the standard TEXTURE_ONLY shader
	return g3d.NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object
}

func NewSceneObject_Globe(rc common.GLRenderingContext) *g3d.SceneObject {
	// Globe model with texture AND normal vectors (for directional lighting)
	geometry := build_globe_geometry(1.0, 64, 32, true)             // create globe geometry with vertex normal vectors
	geometry.BuildDataBuffers(true, false, true)                    // build data buffers for vertices and faces
	material, _ := rc.CreateMaterial("/assets/world.png")           // create material with a texture of world image
	shader := g3d.NewShader_NormalTexture(rc)                       // use the standard NORMAL+TEXTURE shader
	return g3d.NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object
}

func build_globe_geometry(radius float32, wsegs int, hsegs int, use_normals bool) *g3d.Geometry {
	// Globe (sphere) geometry with UV coordinates per vertex (to be used with a texture image)
	//   Note that multiple vertices are assigned to north/south poles, as well as 0/360 longitude.
	//   This approach results in more efficient data buffers than a simple sphere,
	//   since we can build the buffers with single point per vertex, without any repetition.
	geometry := g3d.NewGeometry()
	wnum, hnum := wsegs+1, hsegs+1
	wstep := math.Pi * 2.0 / float32(wsegs)
	hstep := math.Pi / float32(hsegs)
	for i := 0; i < wnum; i++ {
		lon := wstep*float32(i) - math.Pi // longitude (λ) [-180 ~ 180]
		for j := 0; j < hnum; j++ {
			lat := -math.Pi/2.0 + hstep*float32(j) // latitude (φ)
			xyz := GetXYZFromLL(lon, lat, radius)
			geometry.AddVertex(xyz)
			geometry.AddTextureUV([]float32{float32(i) / float32(wsegs), 1.0 - float32(j)/float32(hsegs)})
			if use_normals {
				geometry.AddNormal(c3d.Normalize(xyz))
			}
		}
	}
	for i := 0; i < wnum-1; i++ { // faces on the side
		for j := 0; j < hnum-1; j++ {
			start := uint32((i+0)*hnum + j)
			wnext := uint32((i+1)*hnum + j)
			if spole := (j == 0); spole {
				geometry.AddFace([]uint32{start, wnext + 1, start + 1}) // triangular face for south pole
			} else if npole := (j == hsegs-1); npole {
				geometry.AddFace([]uint32{start, wnext + 0, start + 1}) // triangular face for north pole
			} else {
				geometry.AddFace([]uint32{start, wnext, wnext + 1, start + 1}) // quadratic face otherwise
			}
		}
	}
	return geometry
}

// ----------------------------------------------------------------------------
// GlowRing around the Globe
// ----------------------------------------------------------------------------

func NewSceneObject_GlowRing(rc common.GLRenderingContext) *g3d.SceneObject {
	// GlowRing around the globe, to make the globe stand out against black background.
	// (Note that GlowRing should be rendered in CAMERA space by Renderer)
	if true {
		geometry := build_glowring_geometry(1.0, 1.1, 64)                  // create geometry (a ring around the globe)
		geometry.BuildDataBuffers(true, false, true)                       // build data buffers for vertices and faces
		material, _ := rc.CreateMaterial("GLOW", "#445566")                // texture material for glow effect
		shader := g3d.NewShader_TextureOnly(rc)                            // use the standard TEXTURE_ONLY shader
		scnobj := g3d.NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object
		scnobj.UseBlend = true                                             // default is false
		return scnobj
	} else {
		geometry := build_glowring_geometry(1.0, 1.1, 64)             // create geometry (a ring around the globe)
		geometry.BuildDataBuffersForWireframe()                       // build data buffers for vertices and faces
		shader := g3d.NewShader_ColorOnly(rc)                         // use the standard TEXTURE_ONLY shader
		scnobj := g3d.NewSceneObject(geometry, nil, nil, shader, nil) // set up the scene object
		return scnobj
	}
}

func build_glowring_geometry(in_radius float32, out_radius float32, nsegs int) *g3d.Geometry {
	geometry := g3d.NewGeometry()
	rad := math.Pi * 2 / float64(nsegs)
	for i := 0; i < nsegs; i++ {
		cos, sin := float32(math.Cos(rad*float64(i))), float32(math.Sin(rad*float64(i)))
		geometry.AddVertex([3]float32{in_radius * cos, in_radius * sin, 0})
		geometry.AddVertex([3]float32{out_radius * cos, out_radius * sin, 0})
		geometry.AddTextureUV([]float32{0.0, 0}) // diminishing glow starts
		geometry.AddTextureUV([]float32{1.0, 0}) // diminishing glow ends
		ii, jj := uint32(i), uint32((i+1)%nsegs)
		geometry.AddFace([]uint32{2*ii + 0, 2*jj + 0, 2*jj + 1, 2*ii + 1})
	}
	return geometry
}
