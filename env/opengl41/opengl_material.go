package opengl41

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go4orward/gigl"
)

type OpenGLMaterial struct {
	rc              *OpenGLRenderingContext //
	color           [4][4]float32           // color ([0]:common, [1]:vert, [2]:edge, [3]:face)
	texture         uint32                  // texture
	texture_wh      [2]int                  // texture size
	alphabet_cwh    [2]float32              // character width & height of ALPHABET_STRING
	texture_loading bool                    // true, only if texture is being loaded
	err             error                   //
}

func new_opengl_material(rc *OpenGLRenderingContext, source string, options ...interface{}) (*OpenGLMaterial, error) {
	// THIS CONSTRUCTOR FUNCTION IS NOT MEANT TO BE CALLED DIRECTLY BY USER.
	// IT SHOULD BE CALLED BY 'OpenGLRenderingContext.CreateMaterial()'.
	// 'source' : "#ffffff"
	// 			: "texturepath"
	//			: "GLOW"
	//			: "ALPHABET", foneSize, "#ffffff", outline
	self := OpenGLMaterial{rc: rc, texture: 0, texture_wh: [2]int{0, 0}}
	self.SetDrawModeColor(0, [4]float32{0, 1, 1, 1})
	if source == "" {
		// do nothing
	} else if source == "ALPHABET" { // ALPHABET TEXTURE
		color, fontsize, outline := options[0].(string), options[1].(int), options[2].(bool)
		self.InitializeAlphabetTexture(color, fontsize, outline)
	} else if source == "GLOW" { // 	GLOW TEXTURE
		color := options[0].(string)
		self.InitializeGlowTexture(color)
	} else if source[0] == '#' { // 	RGB COLOR
		rgba := parse_color_from_hex_string(source)
		self.SetDrawModeColor(0, rgba)
	} else { // 						TEXTURE IMAGE PATH
		self.LoadTexture(source)
	}
	return &self, nil
}

func (self *OpenGLMaterial) ShowInfo() {
	colors := ""
	for i := 0; i < len(self.color); i++ {
		c := [4]uint8{uint8(self.color[i][0] * 255), uint8(self.color[i][1] * 255), uint8(self.color[i][2] * 255), uint8(self.color[i][3] * 255)}
		colors += fmt.Sprintf("#%02x%02x%02x%02x ", c[0], c[1], c[2], c[3])
	}
	fmt.Printf("Material with TEXTURE %dx%d and COLOR %s\n", self.texture_wh[0], self.texture_wh[1], colors)
}

// ----------------------------------------------------------------------------
// COLOR
// ----------------------------------------------------------------------------

func (self *OpenGLMaterial) SetColorForDrawMode(draw_mode int, color string) gigl.GLMaterial {
	// 'draw_mode' :  0:common, 1:vertex, 2:edges, 3:faces
	if color != "" {
		return self.SetDrawModeColor(draw_mode, parse_color_from_hex_string(color))
	} else {
		return self
	}
}

func (self *OpenGLMaterial) SetDrawModeColor(draw_mode int, color [4]float32) gigl.GLMaterial {
	switch draw_mode {
	case 1:
		self.color[1] = color // vertex color
	case 2:
		self.color[2] = color // edge color
	case 3:
		self.color[3] = color // face color
	default:
		self.color[0] = color // otherwise
		self.color[1] = color
		self.color[2] = color
		self.color[3] = color
	}
	return self
}

func (self *OpenGLMaterial) GetDrawModeColor(draw_mode int) [4]float32 {
	return self.color[draw_mode]
}

// ----------------------------------------------------------------------------
// TEXTURE
// ----------------------------------------------------------------------------

func (self *OpenGLMaterial) GetTexture() interface{} {
	return self.texture
}

func (self *OpenGLMaterial) GetTextureWH() [2]int {
	return self.texture_wh
}

func (self *OpenGLMaterial) IsTextureReady() bool {
	return (self.texture != 0 && self.texture_wh[0] > 0 && self.texture_wh[1] > 0)
}

func (self *OpenGLMaterial) IsTextureLoading() bool {
	// It may take a long time to load the texture image (ONLY IF it's done asynchronously).
	// That why we have this function defined in GLMaterial interface.
	return self.texture_loading
}

// ----------------------------------------------------------------------------
// Loading Texture Image
// ----------------------------------------------------------------------------

func (self *OpenGLMaterial) LoadTexture(texture_file_path string) *OpenGLMaterial {
	c := self.rc.GetConstants()
	if self.texture == 0 {
		// initialize it with a single CYAN pixel
		var texture uint32
		gl.GenTextures(1, &texture) // uint32, with non-zero values
		gl.BindTexture(c.TEXTURE_2D, texture)
		gl.TexImage2D(c.TEXTURE_2D, 0, int32(c.RGBA), 1, 1, 0, c.RGBA, c.UNSIGNED_BYTE, gl.Ptr([]uint8{0, 255, 255, 255}))
		self.texture = texture
		self.texture_wh = [2]int{1, 1}
		// fmt.Printf("Texture : %v (%T)\n", self.texture, self.texture)
	}
	if texture_file_path != "" {
		// Load texture image asynchronously from a file
		self.texture_loading = true
		// fmt.Printf("Texture started loading %s\n", texture_file_path)
		imgFile, err := os.Open(texture_file_path)
		if err != nil {
			fmt.Printf("WARNING: texture %q not found on disk : %v\n", texture_file_path, err)
			self.err = fmt.Errorf("texture %q not found", texture_file_path)
		} else {
			// img, _, err := image.Decode(imgFile)
			var img image.Image
			switch filepath.Ext(texture_file_path) {
			case ".png", ".PNG":
				img, err = png.Decode(imgFile)
			case ".jpg", ".JPG":
				img, err = jpeg.Decode(imgFile)
			default:
				err = fmt.Errorf("invalid texture file format '%s'", filepath.Ext(texture_file_path))
			}
			imgFile.Close()
			if err != nil {
				fmt.Printf("WARNING: texture %q failed to decode : %v\n", texture_file_path, err)
				self.err = fmt.Errorf("texture %q failed to decode", texture_file_path)
			} else {
				size := img.Bounds().Size()
				// fmt.Printf("Texture image (%dx%d) decoded as %T\n", size.X, size.Y, img)
				var pixbuf []uint8
				switch img.(type) {
				case *image.RGBA: // traditional 32-bit alpha-premultiplied R/G/B/A color
					pixbuf = img.(*image.RGBA).Pix
				case *image.NRGBA: // non-alpha-premultiplied 32-bit R/G/B/A color
					pixbuf = img.(*image.NRGBA).Pix
				default: // unfortunately, we have to convert pixel format
					pixbuf = make([]uint8, size.X*size.Y*4)
					for y := 0; y < size.Y; y++ {
						y_idx := y * size.X * 4
						for x := 0; x < size.X; x++ {
							rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
							idx := y_idx + x*4
							pixbuf[idx+0], pixbuf[idx+1], pixbuf[idx+2], pixbuf[idx+3] = rgba.R, rgba.G, rgba.B, rgba.A
						}
					}
				}
				var texture uint32
				gl.GenTextures(1, &texture) // uint32, with non-zero values
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(c.TEXTURE_2D, texture)
				gl.TexImage2D(c.TEXTURE_2D, 0, int32(c.RGBA), int32(size.X), int32(size.Y), 0, c.RGBA, c.UNSIGNED_BYTE, gl.Ptr(pixbuf))
				if size.X&(size.X-1) == 0 && size.Y&(size.Y-1) == 0 { // POWER-OF-2 width & height
					gl.GenerateMipmap(c.TEXTURE_2D)
				} else { // NON-POWER-OF-2 textures : CLAMP_TO_EDGE & NEAREST/LINEAR only
					gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_S, int32(c.CLAMP_TO_EDGE))
					gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_T, int32(c.CLAMP_TO_EDGE))
					gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MIN_FILTER, int32(c.LINEAR))
					// gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MAG_FILTER, int32(c.LINEAR))
				}
				self.texture = texture
				self.texture_wh = [2]int{size.X, size.Y}
				// fmt.Printf("Texture '%s' %v ready for OpenGL\n", texture_file_path, self.texture_wh)
			}
		}
		self.texture_loading = false
	}
	return self
}

// ----------------------------------------------------------------------------
// Glow Texture
// ----------------------------------------------------------------------------

func (self *OpenGLMaterial) InitializeGlowTexture(color string) {
	rgba := [4]float32{0, 1, 1, 1}
	if len(color) > 0 && color[0] == '#' { // decide RGB color
		rgba = parse_color_from_hex_string(color)
	}
	self.SetDrawModeColor(0, rgba)
	// build pixel buffer for glow effect
	const width, height = 34, 2 // it has to be non-power-of-two texture with gl.NEAREST
	pixbuf := make([]uint8, (width*height)*4)
	for u := 1; u < width-1; u++ { // first (i==0) and last (i==width-1) pixel is ZERO
		ratio := (float32(u-1) / float32(width-2))
		if true { // diminishing glow for the first row (v == 0)  [ 1.0 ~ 0.5 ~ 0.0 ]
			intensity := 1.0 - ratio
			ii := intensity * intensity
			set_pixbuf_with_rgba(pixbuf, (u)*4, uint8(ii*rgba[0]*255), uint8(ii*rgba[1]*255), uint8(ii*rgba[2]*255), uint8(ii*255))
		}
		if true { // glow on both side for the second row (v == 1)  [ 0.0 ~ 1.0 ~ 0.0 ]
			intensity := 1.0 - float32(math.Abs(float64(ratio*2-1)))
			ii := intensity * intensity
			set_pixbuf_with_rgba(pixbuf, (width+u)*4, uint8(ii*rgba[0]*255), uint8(ii*rgba[1]*255), uint8(ii*rgba[2]*255), uint8(ii*255))
		}
	}
	c := self.rc.GetConstants()
	if self.texture == 0 {
		gl.GenTextures(1, &self.texture)
	}
	gl.BindTexture(c.TEXTURE_2D, self.texture)
	gl.TexImage2D(c.TEXTURE_2D, 0, int32(c.RGBA), width, height, 0, c.RGBA, c.UNSIGNED_BYTE, gl.Ptr(pixbuf))
	gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_S, int32(c.CLAMP_TO_EDGE))
	gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_WRAP_T, int32(c.CLAMP_TO_EDGE))
	gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MIN_FILTER, int32(c.LINEAR))
	// gl.TexParameteri(c.TEXTURE_2D, c.TEXTURE_MAG_FILTER, int32(c.LINEAR))
	self.texture_wh = [2]int{width, height} // CLAMP_TO_EDGE & NEAREST(not LINEAR) for NON-POWER-OF-2 textures
}

func set_pixbuf_with_rgba(pbuffer []uint8, idx int, R uint8, G uint8, B uint8, A uint8) {
	pbuffer[idx+0] = R
	pbuffer[idx+1] = G
	pbuffer[idx+2] = B
	pbuffer[idx+3] = A
}

// ----------------------------------------------------------------------------
// Alphabet Texture  (for labels)
// ----------------------------------------------------------------------------

const _ALPHABET_STRING = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~?°"

func (self *OpenGLMaterial) InitializeAlphabetTexture(color string, fontsize int, outlined bool) {
	// TODO
}

func (self *OpenGLMaterial) GetAlaphabetLength() int {
	return len([]rune(_ALPHABET_STRING)) // total rune length of _ALPHABET_STRING
}

func (self *OpenGLMaterial) GetAlaphabetCharacterWH(scale float32) [2]float32 {
	return [2]float32{self.alphabet_cwh[0] * scale, self.alphabet_cwh[1] * scale}
}

func (self *OpenGLMaterial) GetAlaphabetCharacterIndex(c rune) int {
	if c >= 32 && c < 127 { // 32: <SPACE>, 127: <DEL>
		return int(c) - 32
	} else {
		alphabet_rune := []rune(_ALPHABET_STRING)
		for i := 127 - 32; i < len(alphabet_rune); i++ {
			if c == alphabet_rune[i] {
				return i
			}
		}
		return (127 - 32) // '?'
	}
}

// ----------------------------------------------------------------------------
// private function
// ----------------------------------------------------------------------------

func parse_color_from_hex_string(s string) [4]float32 {
	c := [4]uint8{0, 0, 0, 255}
	if len(s) == 0 {
	} else if s[0] == '#' {
		switch len(s) {
		case 9:
			fmt.Sscanf(s, "#%02x%02x%02x%02x", &c[0], &c[1], &c[2], &c[3])
		case 7:
			fmt.Sscanf(s, "#%02x%02x%02x", &c[0], &c[1], &c[2])
		case 5:
			fmt.Sscanf(s, "#%1x%1x%1x%1x", &c[0], &c[1], &c[2], &c[3])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
			c[3] *= 17
		case 4:
			fmt.Sscanf(s, "#%1x%1x%1x", &c[0], &c[1], &c[2])
			c[0] *= 17
			c[1] *= 17
			c[2] *= 17
		default:
		}
	} else {
	}
	return [4]float32{float32(c[0]) / 255, float32(c[1]) / 255, float32(c[2]) / 255, float32(c[3]) / 255}
}
