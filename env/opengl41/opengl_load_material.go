package opengl41

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/g2d"
)

func load_material(rc *OpenGLRenderingContext, material gigl.GLMaterial) error {
	switch material.(type) {
	case *g2d.MaterialColors:
		prepare_material_colors(rc, material.(*g2d.MaterialColors))
	case *g2d.MaterialTexture:
		prepare_material_texture(rc, material.(*g2d.MaterialTexture))
	case *g2d.MaterialGlowTexture:
		prepare_material_glow_texture(rc, material.(*g2d.MaterialGlowTexture))
	case *g2d.MaterialAlphabetTexture:
		prepare_material_alphabet_texture(rc, material.(*g2d.MaterialAlphabetTexture))
	}
	return nil
}

// ----------------------------------------------------------------------------
// Material Colors
// ----------------------------------------------------------------------------

func prepare_material_colors(rc *OpenGLRenderingContext, mcolor *g2d.MaterialColors) {
	// Nothing to be done
}

// ----------------------------------------------------------------------------
// Material Texture
// ----------------------------------------------------------------------------

func prepare_material_texture(rc *OpenGLRenderingContext, mtex *g2d.MaterialTexture) {
	c := rc.GetConstants()
	if !mtex.IsTextureReady() && !mtex.IsTextureLoading() {
		if true {
			// set up a temporary texture (single pixel with CYAN colar)
			var texture uint32
			gl.GenTextures(1, &texture) // uint32, with non-zero values
			gl.BindTexture(c.TEXTURE_2D, texture)
			gl.TexImage2D(c.TEXTURE_2D, 0, int32(c.RGBA), 1, 1, 0, c.RGBA, c.UNSIGNED_BYTE, gl.Ptr([]uint8{0, 255, 255, 255}))
			mtex.SetTexture(texture)
			// fmt.Printf("Texture : %v (%T)\n", self.texture, self.texture)
		}
		// get the pixel buffer, and the width & height of the texture
		mtex.LoadTextureFromLocalFile(func(pixbuf []uint8, wh [2]int) {
			var texture uint32
			gl.GenTextures(1, &texture)   // If gl functions run on different threads other than main, it fails.
			gl.ActiveTexture(gl.TEXTURE0) // Therefore, LoadTextureFromLocalFile() cannot be asynchronous.
			gl.BindTexture(c.TEXTURE_2D, texture)
			gl.TexImage2D(c.TEXTURE_2D, 0, int32(c.RGBA), int32(wh[0]), int32(wh[1]), 0, c.RGBA, c.UNSIGNED_BYTE, gl.Ptr(pixbuf))
			if wh[0]&(wh[0]-1) == 0 && wh[1]&(wh[1]-1) == 0 { // POWER-OF-2 width & height
				gl.GenerateMipmap(c.TEXTURE_2D)
			} else { // NON-POWER-OF-2 textures : CLAMP_TO_EDGE & NEAREST/LINEAR only
				gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_S, int32(c.CLAMP_TO_EDGE))
				gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_T, int32(c.CLAMP_TO_EDGE))
				gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MIN_FILTER, int32(c.LINEAR))
				// gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MAG_FILTER, int32(c.LINEAR))
			}
			mtex.SetTexture(texture)
		})
	}
}

// ----------------------------------------------------------------------------
// Material Glow Texture
// ----------------------------------------------------------------------------

func prepare_material_glow_texture(rc *OpenGLRenderingContext, mtex *g2d.MaterialGlowTexture) {
	c := rc.GetConstants()
	if !mtex.IsTextureReady() {
		// get the pixel buffer, and the width & height of the texture
		mtex.LoadTexture(func(pixbuf []uint8, wh [2]int) {
			var texture uint32
			gl.GenTextures(1, &texture)
			gl.BindTexture(c.TEXTURE_2D, texture)
			gl.TexImage2D(c.TEXTURE_2D, 0, int32(c.RGBA), int32(wh[0]), int32(wh[1]), 0, c.RGBA, c.UNSIGNED_BYTE, gl.Ptr(pixbuf))
			gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_S, int32(c.CLAMP_TO_EDGE))
			gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_T, int32(c.CLAMP_TO_EDGE))
			gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MIN_FILTER, int32(c.LINEAR))
			// gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MAG_FILTER, int32(c.LINEAR))
			mtex.SetTexture(texture)
		})
	}
}

// ----------------------------------------------------------------------------
// Material Alphabet Texture
// ----------------------------------------------------------------------------

func prepare_material_alphabet_texture(rc *OpenGLRenderingContext, mtex *g2d.MaterialAlphabetTexture) {
	// TODO(go4orward)
}
