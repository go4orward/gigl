package earth

import (
	"math"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/common"
	"github.com/go4orward/gigl/g2d"
	"github.com/go4orward/gigl/g3d"
)

type WorldGlobe struct {
	bkgcolor    [3]float32       // background color of the globe
	GSphere     *g3d.SceneObject // globe sphere with texture & vertex normals
	GlowRing    *g3d.SceneObject // glow ring around the globe
	modelmatrix common.Matrix4   // Model matrix of the globe & its layers
}

func NewWorldGlobe(rc gigl.GLRenderingContext, bkg_color string, world_img_filepath string) *WorldGlobe {
	self := WorldGlobe{}
	self.SetBkgColor(bkg_color)
	// Globe
	if true { // with texture AND normal vectors (for directional lighting)
		use_vertex_normal := true
		geometry := build_globe_geometry(1.0, 64, 32, use_vertex_normal)        // create globe geometry with vertex normal vectors
		geometry.BuildDataBuffers(true, false, true)                            // build data buffers for vertices and faces
		material := g2d.NewMaterialTexture(world_img_filepath)                  // create material with a texture of world image
		shader := g3d.NewShader_NormalTexture(rc)                               // use the standard NORMAL+TEXTURE shader
		self.GSphere = g3d.NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object
	} else { // with texture UV coordinates only (without normal vectors and directional lighting)
		use_vertex_normal := false
		geometry := build_globe_geometry(1.0, 64, 32, use_vertex_normal)        // create globe geometry with texture UVs only
		geometry.BuildDataBuffers(true, false, true)                            // build data buffers for vertices and faces
		material := g2d.NewMaterialTexture(world_img_filepath)                  // create material with a texture of world image
		shader := g3d.NewShader_TextureOnly(rc)                                 // use the standard TEXTURE_ONLY shader
		self.GSphere = g3d.NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object
	}
	if common.Logger.IsLogging(common.LogLevelTrace) {
		common.Logger.Trace("WorldGlobe \n%s", self.GSphere.Summary())
	}
	// GlowRing around the globe (to make the globe stand out against black background)
	//   (Note that GlowRing should be rendered in CAMERA space by Renderer)
	if true {
		geometry := build_glowring_geometry(1.0, 1.1, 64)                  // create geometry (a ring around the globe)
		geometry.BuildDataBuffers(true, false, true)                       // build data buffers for vertices and faces
		material := g2d.NewMaterialGlowTexture("#445566")                  // texture material for glow effect
		rc.LoadMaterial(material)                                          // load the texture
		shader := g3d.NewShader_TextureOnly(rc)                            // use the standard TEXTURE_ONLY shader
		scnobj := g3d.NewSceneObject(geometry, material, nil, nil, shader) // set up the scene object
		scnobj.UseBlend = true                                             // default is false
		self.GlowRing = scnobj
	} else {
		geometry := build_glowring_geometry(1.0, 1.1, 64)             // create geometry (a ring around the globe)
		geometry.BuildDataBuffersForWireframe()                       // build data buffers for vertices and faces
		shader := g3d.NewShader_ColorOnly(rc)                         // use the standard TEXTURE_ONLY shader
		scnobj := g3d.NewSceneObject(geometry, nil, nil, shader, nil) // set up the scene object
		self.GlowRing = scnobj
	}
	self.modelmatrix.SetIdentity() // initialize as Identity matrix
	return &self
}

// ----------------------------------------------------------------------------
// Background Color
// ----------------------------------------------------------------------------

func (self *WorldGlobe) SetBkgColor(color string) *WorldGlobe {
	rgba := common.RGBAFromHexString(color)
	self.bkgcolor = [3]float32{rgba[0], rgba[1], rgba[2]}
	return self
}

func (self *WorldGlobe) GetBkgColor() [3]float32 {
	return self.bkgcolor
}

// ----------------------------------------------------------------------------
// Translation, Rotation, Scaling (by manipulating MODEL matrix)
// ----------------------------------------------------------------------------

func (self *WorldGlobe) SetTransformation(txyz [3]float32, axis [3]float32, angle_in_degree float32, scale float32) *WorldGlobe {
	translation := common.NewMatrix4().SetTranslation(txyz[0], txyz[1], txyz[2])
	rotation := common.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	scaling := common.NewMatrix4().SetScaling(scale, scale, scale)
	self.modelmatrix.SetMultiplyMatrices(translation, rotation, scaling)
	return self
}

func (self *WorldGlobe) Translate(tx float32, ty float32, tz float32) *WorldGlobe {
	translation := common.NewMatrix4().SetTranslation(tx, ty, tz)
	self.modelmatrix.SetMultiplyMatrices(translation, &self.modelmatrix)
	return self
}

func (self *WorldGlobe) Rotate(axis [3]float32, angle_in_degree float32) *WorldGlobe {
	rotation := common.NewMatrix4().SetRotationByAxis(axis, angle_in_degree)
	self.modelmatrix.SetMultiplyMatrices(rotation, &self.modelmatrix)
	return self
}

func (self *WorldGlobe) Scale(scale float32) *WorldGlobe {
	scaling := common.NewMatrix4().SetScaling(scale, scale, scale)
	self.modelmatrix.SetMultiplyMatrices(scaling, &self.modelmatrix)
	return self
}

// ----------------------------------------------------------------------------
// WorldGlobe
// ----------------------------------------------------------------------------

func build_globe_geometry(radius float32, wsegs int, hsegs int, use_normals bool) *g3d.Geometry {
	// WorldGlobe (sphere) geometry with UV coordinates per vertex (to be used with a texture image)
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
			geometry.AddVertex(*xyz)
			geometry.AddTextureUV([]float32{float32(i) / float32(wsegs), 1.0 - float32(j)/float32(hsegs)})
			if use_normals {
				geometry.AddNormal(*xyz.Normalize())
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
// GlowRing around the WorldGlobe
// ----------------------------------------------------------------------------

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
