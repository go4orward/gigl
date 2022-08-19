package webgl10

import (
	"fmt"
	"math"
	"syscall/js"

	"github.com/go4orward/gigl"
	"github.com/go4orward/gigl/g2d"
)

func load_material(rc *WebGLRenderingContext, material gigl.GLMaterial) error {
	// Load material data from local file or remote server
	context, c := rc.context, rc.GetConstants()
	switch material.(type) {
	case *g2d.MaterialColors:
		// DO NOTHING
	case *g2d.MaterialTexture:
		mtex := material.(*g2d.MaterialTexture)
		if !mtex.IsReady() && !mtex.IsLoaded() && !mtex.IsLoading() {
			if false {
				// set up a temporary texture (single pixel with CYAN colar)
				mtex.SetTexture(context.Call("createTexture"))
				// js_texture_unit := js.ValueOf(js.ValueOf(self.constants.TEXTURE0 + uint32(texture_unit)))
				context.Call("activeTexture", js.ValueOf(c.TEXTURE0))
				context.Call("bindTexture", js.ValueOf(c.TEXTURE_2D), mtex.GetTexture())
				// context.TexImage2DFromPixelBuffer(c.TEXTURE_2D, 0, c.RGBA, 1, 1, 0, c.RGBA, c.UNSIGNED_BYTE, []uint8{0, 255, 255, 255})
				js_buffer := rc.ConvertGoSliceToJsTypedArray([]uint8{0, 255, 255, 255})
				context.Call("texImage2D", js.ValueOf(c.TEXTURE_2D), 0, js.ValueOf(c.RGBA), 1, 1, 0, js.ValueOf(c.RGBA), js.ValueOf(c.UNSIGNED_BYTE), js_buffer)
			}
			// get the pixel buffer, and the width & height of the texture
			mtex.LoadTextureFromRemoteServer()
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

func setup_material(rc *WebGLRenderingContext, material gigl.GLMaterial) error {
	// Setup material using pre-loaded material data
	context, c := rc.context, rc.GetConstants()
	switch material.(type) {
	case *g2d.MaterialColors:
		// DO NOTHING
	case *g2d.MaterialTexture:
		mtex := material.(*g2d.MaterialTexture)
		if !mtex.IsReady() && mtex.IsLoaded() {
			pixbuf, wh := mtex.GetTexturePixbuf(), mtex.GetTextureWH()
			mtex.SetTexture(context.Call("createTexture"))
			context.Call("bindTexture", js.ValueOf(c.TEXTURE_2D), mtex.GetTexture())
			// rc.GLTexImage2DFromPixelBuffer(c.TEXTURE_2D, 0, c.RGBA, wh[0], wh[1], 0, c.RGBA, c.UNSIGNED_BYTE, pixbuf)
			js_buffer := rc.ConvertGoSliceToJsTypedArray(pixbuf)
			context.Call("texImage2D", js.ValueOf(c.TEXTURE_2D), 0, js.ValueOf(c.RGBA), wh[0], wh[1], 0, js.ValueOf(c.RGBA), js.ValueOf(c.UNSIGNED_BYTE), js_buffer)
			if wh[0]&(wh[0]-1) == 0 && wh[1]&(wh[1]-1) == 0 { // POWER-OF-2 width & height
				context.Call("generateMipmap", js.ValueOf(c.TEXTURE_2D))
			} else { // NON-POWER-OF-2 textures : CLAMP_TO_EDGE & NEAREST/LINEAR only
				context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_WRAP_S), js.ValueOf(c.CLAMP_TO_EDGE))
				context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_WRAP_T), js.ValueOf(c.CLAMP_TO_EDGE))
				context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_MIN_FILTER), js.ValueOf(c.LINEAR))
			}
		}
	case *g2d.MaterialGlowTexture:
		mtex := material.(*g2d.MaterialGlowTexture)
		if !mtex.IsReady() && mtex.IsLoaded() {
			pixbuf, wh := mtex.GetTexturePixbuf(), mtex.GetTextureWH()
			mtex.SetTexture(context.Call("createTexture"))
			context.Call("bindTexture", js.ValueOf(c.TEXTURE_2D), mtex.GetTexture())
			js_buffer := rc.ConvertGoSliceToJsTypedArray(pixbuf)
			context.Call("texImage2D", js.ValueOf(c.TEXTURE_2D), 0, js.ValueOf(c.RGBA), wh[0], wh[1], 0, js.ValueOf(c.RGBA), js.ValueOf(c.UNSIGNED_BYTE), js_buffer)
			context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_WRAP_S), js.ValueOf(c.CLAMP_TO_EDGE))
			context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_WRAP_T), js.ValueOf(c.CLAMP_TO_EDGE))
			context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_MIN_FILTER), js.ValueOf(c.LINEAR))
		}
	case *g2d.MaterialAlphabetTexture:
	}
	return nil
}

// ----------------------------------------------------------------------------
// Alphabet Texture  (for labels)
// ----------------------------------------------------------------------------

func prepare_material_alphabet_texture(rc *WebGLRenderingContext, mab *g2d.MaterialAlphabetTexture) {
	// 'fontsize' : 12=>(7.2x12.6), 16=>(9.6x16.8), 20=>(12x21), 24=>(14x25), 30=>(18x31), 40=>(24x42)
	font_style := fmt.Sprintf("%dpx %s", mab.GetFontSize(), mab.GetFontFamily())
	txtctx := js.Global().Get("document").Call("createElement", "canvas").Call("getContext", "2d")
	txtctx.Set("font", font_style) // need to be set, before measuring text size
	cwidth := float32(txtctx.Call("measureText", "M").Get("width").Float())
	cheight := float32(mab.GetFontSize()) * 1.05 // we need some more margin below the text
	twidth := int(math.Floor(txtctx.Call("measureText", mab.GetAlphabetString()).Get("width").Float()))
	theight := int(cheight) // instead of int(cwidth*2)
	// common.Logger.Trace("Character: %v %v  Texture: %v %v\n", cwidth, cheight, twidth, theight)
	txtctx.Get("canvas").Set("width", twidth)
	txtctx.Get("canvas").Set("height", theight)
	txtctx.Call("clearRect", 0, 0, twidth, theight)
	txtctx.Set("font", font_style)    // need to be set again!
	txtctx.Set("textAlign", "start")  // start (default), end, left, right, center
	txtctx.Set("textBaseline", "top") // top, hanging, middle, alphabetic (default), ideographic, bottom
	if mab.GetFontOutlined() {
		txtctx.Set("strokeStyle", "#000000")                     // BLACK outline
		txtctx.Set("lineWidth", 2.5)                             // text stroke width
		txtctx.Call("strokeText", mab.GetAlphabetString(), 0, 0) // draw the alphabet string for outline
	}
	txtctx.Set("fillStyle", mab.GetFontColor())            // interior (Note that WHITE can be multiplied with other colors later)
	txtctx.Call("fillText", mab.GetAlphabetString(), 0, 0) // draw the alphabet string
	context, c := rc.context, rc.GetConstants()
	mab.SetTexture(context.Call("createTexture"))
	context.Call("bindTexture", js.ValueOf(c.TEXTURE_2D), mab.GetTexture())
	// mab.rc.GLTexImage2DFromImgObject(c.TEXTURE_2D, 0, c.RGBA, c.RGBA, c.UNSIGNED_BYTE, txtctx.Get("canvas"))
	context.Call("texImage2D", js.ValueOf(c.TEXTURE_2D), 0, js.ValueOf(c.RGBA), js.ValueOf(c.RGBA), js.ValueOf(c.UNSIGNED_BYTE), txtctx.Get("canvas"))
	context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_WRAP_S), js.ValueOf(c.CLAMP_TO_EDGE))
	context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_WRAP_T), js.ValueOf(c.CLAMP_TO_EDGE))
	context.Call("texParameteri", js.ValueOf(c.TEXTURE_2D), js.ValueOf(c.TEXTURE_MIN_FILTER), js.ValueOf(c.LINEAR))
	mab.SetTextureWH([2]int{twidth, theight})
	mab.SetAlphabetWH([2]float32{cwidth, cheight})
}
