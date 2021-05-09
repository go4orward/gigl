package g3d

import (
	"fmt"

	"github.com/go4orward/gigl/common"
)

type Scene struct {
	bkgcolor [3]float32     // background color of the scene
	objects  []*SceneObject // SceneObjects in the scene
	overlays []Overlay      // list of Overlay (interface) layers
}

func NewScene(bkg_color string) *Scene {
	var scene Scene
	scene.SetBkgColor(bkg_color)
	scene.objects = make([]*SceneObject, 0)
	scene.overlays = make([]Overlay, 0)
	return &scene
}

func (self *Scene) String() string {
	return fmt.Sprintf("2DScene{objects:%d overlays:%d}\n", len(self.objects), len(self.overlays))
}

// ----------------------------------------------------------------------------
// Background Color
// ----------------------------------------------------------------------------

func (self *Scene) SetBkgColor(color string) *Scene {
	rgba := common.ParseHexColor(color)
	self.bkgcolor = [3]float32{rgba[0], rgba[1], rgba[2]}
	return self
}

func (self *Scene) GetBkgColor() [3]float32 {
	return self.bkgcolor
}

// ----------------------------------------------------------------------------
// Handling SceneObject
// ----------------------------------------------------------------------------

func (self *Scene) Add(scnobj ...*SceneObject) *Scene {
	for i := 0; i < len(scnobj); i++ {
		self.objects = append(self.objects, scnobj[i])
	}
	return self
}

func (self *Scene) Get(indices ...int) *SceneObject {
	// Find a SceneObject using the list of indices
	// (multiple indices refers to children[i] of SceneObject)
	scene_object_list := self.objects
	for i := 0; i < len(indices); i++ {
		index := indices[i]
		if index < 0 || index >= len(scene_object_list) {
			return nil
		} else if i == len(indices)-1 {
			scene_object := scene_object_list[index]
			return scene_object
		} else {
			scene_object := scene_object_list[index]
			scene_object_list = scene_object.children
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
// Managing OverlayLayers
// ----------------------------------------------------------------------------

func (self *Scene) AddOverlay(overlay ...Overlay) *Scene {
	for i := 0; i < len(overlay); i++ {
		self.overlays = append(self.overlays, overlay[i])
	}
	return self
}
