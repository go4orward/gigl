package g3d

// func (self *Scene) loadPendulum() *Scene {
// 	pendulum := NewSceneObject(NewGeometry().LoadCylinder(10, 38, 1, false).Translate(0, 0, +0.5), "FFFFFF", "base")
// 	pillar := NewSceneObject(NewGeometry().LoadCube(1, 1, 10, false).Translate(0, 9, 6), "FFFFFF", "pillar")
// 	arm := NewSceneObject(NewGeometry().LoadCylinder(0.3, 12, 10, false).Rotate(1, 0, 0, 90).Translate(0, 4.5, 0), "FFFFFF", "arm")
// 	rope := NewSceneObject(NewGeometry().LoadCylinder(0.1, 12, 7, false).Translate(0, 0, -3.5), "FFFFFF", "rope")
// 	ball := NewSceneObject(NewGeometry().LoadSphere(1, 36, 18, false).Translate(0, 0, -7), "FFFFFF", "ball")
// 	pendulum.AddChild(pillar).AddChild(arm).Translate(0, 0, 10).AddChild(rope.AddChild(ball))
// 	self.Add(pendulum)
// 	return self
// }
