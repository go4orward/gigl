package g2d

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go4orward/gigl/common"
)

type MaterialTexture struct {
	image_filepath  string     //
	pixbuf          []uint8    //
	texture         any        // texture (js.Value for WebGL, uint32 for OpenGL)
	texture_wh      [2]int     // texture size
	texture_rgb     [3]float32 // extra RGB color to be multiplied with the texture
	texture_loading bool       // true, only if texture is being loaded
	err             error      //
}

func NewMaterialTexture(filepath string, color ...string) *MaterialTexture {
	mtex := MaterialTexture{image_filepath: filepath}
	if len(color) > 0 && color[0][0] == '#' {
		mtex.texture_rgb = common.RGBFromHexString(color[0])
	}
	// NOTE THAT ACTUAL PREPARATION OF THIS MATERIAL IS IMPLEMENTED FOR EACH ENVIRONMENT,
	//   IT IS AUTOMATICALLY CALLED BY Renderer USING GLRenderingContext.LoadMaterial(material).
	return &mtex
}

func (self *MaterialTexture) MaterialSummary() string {
	if self.IsReady() || self.IsLoaded() {
		return fmt.Sprintf("MaterialTexture %dx%d %q", self.texture_wh[0], self.texture_wh[1], self.image_filepath)
	} else if self.IsLoading() {
		return fmt.Sprintf("MaterialTexture loading %q", self.image_filepath)
	} else {
		return fmt.Sprintf("MaterialTexture %q", self.image_filepath)
	}
}

// ----------------------------------------------------------------------------
// TEXTURE
// ----------------------------------------------------------------------------

func (self *MaterialTexture) GetTexturePixbuf() []uint8 {
	return self.pixbuf
}

func (self *MaterialTexture) GetTexture() any {
	return self.texture
}

func (self *MaterialTexture) SetTexture(texture any) {
	self.texture = texture
}

func (self *MaterialTexture) GetTextureWH() [2]int {
	return self.texture_wh
}

func (self *MaterialTexture) GetTextureRGB() [3]float32 {
	return self.texture_rgb
}

func (self *MaterialTexture) SetTextureRGB(color any) {
	switch color.(type) {
	case string:
		self.texture_rgb = common.RGBFromHexString(color.(string))
	case [3]float32:
		self.texture_rgb = color.([3]float32)
	case []float32:
		c := color.([]float32)
		self.texture_rgb = [3]float32{c[0], c[1], c[2]}
	}
}

func (self *MaterialTexture) IsReady() bool {
	// Texture was successfully set up, and it's ready for rendering.
	return (self.texture != nil && self.texture_wh[0] > 0 && self.texture_wh[1] > 0 && self.err == nil)
}

func (self *MaterialTexture) IsLoaded() bool {
	// Texture was successfully loaded, and it needs to be set up by main thread.
	return (self.pixbuf != nil && self.texture_wh[0] > 0 && self.texture_wh[1] > 0 && self.err == nil)
}

func (self *MaterialTexture) IsLoading() bool {
	// Texture is being loaded asynchronously by non-main thread (using Go function).
	return self.texture_loading
}

// ----------------------------------------------------------------------------
// Loading Texture Image
// ----------------------------------------------------------------------------

func (self *MaterialTexture) LoadTextureFromLocalFile() {
	if self.image_filepath != "" && self.err == nil {
		// Load texture image asynchronously from a file
		self.texture = nil
		self.pixbuf = nil
		self.texture_wh = [2]int{0, 0}
		self.texture_loading = true
		go func() { // OpenGL's gl.GenTextures() fails on different threads other than main thread.
			// common.Logger.Trace("Texture started loading %s\n", self.image_filepath)
			defer func() { self.texture_loading = false }()
			imgFile, err := os.Open(self.image_filepath)
			if err != nil {
				self.err = fmt.Errorf("texture %q not found", self.image_filepath)
				common.Logger.Error("%v\n", self.err)
			} else {
				defer imgFile.Close()
				bytes, _ := ioutil.ReadAll(imgFile)
				pixbuf, wh, err := self.decode_pixels_from_image_bytes(bytes, filepath.Ext(self.image_filepath))
				if err != nil {
					self.err = fmt.Errorf("texture %q failed to decode (%v)", self.image_filepath, err)
					common.Logger.Error("%v\n", self.err)
				} else {
					self.pixbuf = pixbuf
					self.texture_wh = wh
					// common.Logger.Trace("Texture '%s' %v ready for OpenGL\n", self.image_filepath, self.texture_wh)
				}
			}
		}()
	}
}

func (self *MaterialTexture) LoadTextureFromRemoteServer() {
	if self.image_filepath != "" && self.err == nil {
		self.texture = nil
		self.pixbuf = nil
		self.texture_wh = [2]int{0, 0}
		request_url := self.image_filepath
		self.texture_loading = true
		go func() {
			// common.Logger.Trace("Texture started GET %q\n", request_url)
			defer func() { self.texture_loading = false }()
			resp, err := http.Get(request_url)
			if err != nil {
				self.err = fmt.Errorf("Failed to GET %q (%v)\n", request_url, err)
				common.Logger.Error("%v\n", self.err)
				return
			}
			defer resp.Body.Close()
			response_body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				self.err = fmt.Errorf("Failed to GET %q (%v)\n", request_url, err)
				common.Logger.Error("%v\n", self.err)
			} else if resp.StatusCode < 200 || resp.StatusCode > 299 { // response with error message
				self.err = fmt.Errorf("Failed to GET %s : (%d) %s\n", request_url, resp.StatusCode, string(response_body))
				common.Logger.Error("%v\n", self.err)
			} else { // successful response with texture image
				pixbuf, wh, err := self.decode_pixels_from_image_bytes(response_body, filepath.Ext(request_url))
				// common.Logger.Error("Texture '%s' %v ready for OpenGL\n", self.image_filepath, self.texture_wh)
				if err != nil {
					self.err = fmt.Errorf("Failed to decode %q (err:%v)\n", request_url, err.Error())
					common.Logger.Error("%v\n", self.err)
				} else {
					self.pixbuf = pixbuf
					self.texture_wh = wh
					// log.Printf("Texture (%dx%d) loaded for WebGL from server %q (ready:%v)\n", wh[0], wh[1], request_url, self.IsTextureReady())
				}
			}
		}()
	}
}

func (self *MaterialTexture) decode_pixels_from_image_bytes(imgbytes []byte, ext string) ([]uint8, [2]int, error) {
	var img image.Image
	var err error
	switch ext {
	case ".png", ".PNG":
		img, err = png.Decode(bytes.NewBuffer(imgbytes))
	case ".jpg", ".JPG":
		img, err = jpeg.Decode(bytes.NewBuffer(imgbytes))
	default:
		return nil, [2]int{}, fmt.Errorf("invalid texture image format %q", ext)
	}
	if err != nil {
		return nil, [2]int{}, err
	}
	size := img.Bounds().Size()
	// common.Logger.Trace("Texture image (%dx%d) decoded as %T\n", size.X, size.Y, img)
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
				pixbuf[idx+0] = rgba.R
				pixbuf[idx+1] = rgba.G
				pixbuf[idx+2] = rgba.B
				pixbuf[idx+3] = rgba.A

			}
		}
	}
	return pixbuf, [2]int{size.X, size.Y}, nil
}
