package opengl41

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/g2d"
)

func load_material(rc *OpenGLRenderingContext, material gigl.GLMaterial) error {
	c := rc.GetConstants()
	switch material.(type) {
	case *g2d.MaterialColors:
		// DO NOTHING
	case *g2d.MaterialTexture:
		mtex := material.(*g2d.MaterialTexture)
		if !mtex.IsReady() && !mtex.IsLoaded() && !mtex.IsLoading() {
			if false {
				// set up a temporary texture (single pixel with CYAN colar)
				var texture uint32
				gl.GenTextures(1, &texture) // uint32, with non-zero values
				gl.BindTexture(c.TEXTURE_2D, texture)
				gl.TexImage2D(c.TEXTURE_2D, 0, int32(c.RGBA), 1, 1, 0, c.RGBA, c.UNSIGNED_BYTE, gl.Ptr([]uint8{0, 255, 255, 255}))
				mtex.SetTexture(texture)
				// common.Logger.Trace("Texture : %v (%T)\n", self.texture, self.texture)
			}
			// get the pixel buffer, and the width & height of the texture
			mtex.LoadTextureFromLocalFile()
		}
	case *g2d.MaterialGlowTexture:
		mtex := material.(*g2d.MaterialGlowTexture)
		if !mtex.IsReady() && !mtex.IsLoaded() {
			// get the pixel buffer, and the width & height of the texture
			mtex.LoadGlowTexture()
		}
	case *g2d.MaterialAlphabetTexture:
		prepare_material_alphabet_texture(rc, material.(*g2d.MaterialAlphabetTexture))
	}
	return nil
}

func setup_material(rc *OpenGLRenderingContext, material gigl.GLMaterial) error {
	c := rc.GetConstants()
	switch material.(type) {
	case *g2d.MaterialColors:
		// DO NOTHING
	case *g2d.MaterialTexture:
		mtex := material.(*g2d.MaterialTexture)
		if !mtex.IsReady() && mtex.IsLoaded() {
			pixbuf, wh := mtex.GetTexturePixbuf(), mtex.GetTextureWH()
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
		}
	case *g2d.MaterialGlowTexture:
		mtex := material.(*g2d.MaterialGlowTexture)
		if !mtex.IsReady() && mtex.IsLoaded() {
			pixbuf, wh := mtex.GetTexturePixbuf(), mtex.GetTextureWH()
			var texture uint32
			gl.GenTextures(1, &texture)
			gl.BindTexture(c.TEXTURE_2D, texture)
			gl.TexImage2D(c.TEXTURE_2D, 0, int32(c.RGBA), int32(wh[0]), int32(wh[1]), 0, c.RGBA, c.UNSIGNED_BYTE, gl.Ptr(pixbuf))
			gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_S, int32(c.CLAMP_TO_EDGE))
			gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_T, int32(c.CLAMP_TO_EDGE))
			gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MIN_FILTER, int32(c.LINEAR))
			// gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MAG_FILTER, int32(c.LINEAR))
			mtex.SetTexture(texture)
		}
	case *g2d.MaterialAlphabetTexture:
		prepare_material_alphabet_texture(rc, material.(*g2d.MaterialAlphabetTexture))
	}
	return nil
}

// ----------------------------------------------------------------------------
// Material Alphabet Texture
// ----------------------------------------------------------------------------

func prepare_material_alphabet_texture(rc *OpenGLRenderingContext, mtex *g2d.MaterialAlphabetTexture) {
	// TODO(go4orward)
}
