package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/misc/particles"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/RangerGo/engine/rendering"
)

type gameLayer struct {
	nodes.Node

	dragSquare *draggableSquare

	// Particles
	particleSystem    api.IParticleSystem
	autoTriggerEnable bool
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	if err := o.Build(world); err != nil {
		return nil, err
	}
	return o, nil
}

func (g *gameLayer) Build(world api.IWorld) error {
	g.Node.Build(world)

	var err error

	dvr := world.Properties().Window.DeviceRes

	// ---------------------------------------------------------
	shline, err := custom.NewStaticHLineNode("HLine", world, g)
	if err != nil {
		return err
	}
	shline.SetScale(float32(dvr.Width))
	ghl := shline.(*custom.StaticHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	svline, err := custom.NewStaticVLineNode("VLine", world, g)
	if err != nil {
		return err
	}
	svline.SetScale(float32(dvr.Width))
	gvl := svline.(*custom.StaticVLineNode)
	gvl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	tri, err := custom.NewStaticTriangleNode("FilledTri", true, true, world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100)
	tri.SetPosition(150.0, 0.0)
	gtr := tri.(*custom.StaticTriangleNode)
	gtr.SetColor(color.NewPaletteInt64(color.PanSkin))

	// ---------------------------------------------------------
	otri, err := custom.NewStaticTriangleNode("OutlineTri", true, false, world, g)
	if err != nil {
		return err
	}
	otri.SetScale(100)
	otri.SetPosition(150.0, 0.0)
	gotr := otri.(*custom.StaticTriangleNode)
	gotr.SetColor(color.NewPaletteInt64(color.White))

	// Particle system
	activator := particles.NewActivator360()
	g.particleSystem = particles.NewParticleSystem(activator)
	g.particleSystem.Activate(true)
	g.particleSystem.SetAutoTrigger(false)

	// Now populate the system
	for i := 0; i < 50; i++ {
		ptri, err := newParticleTriangleNode("FilledTri", true, true, world, g)
		if err != nil {
			return err
		}
		gptri := ptri.(*ParticleTriangleNode)
		gptri.SetColor(color.NewPaletteInt64(rendering.Black).Array())
		gptri.SetVisible(false)
		gptri.SetScale(10.0)
		p := particles.NewParticle(ptri)
		g.particleSystem.AddParticle(p)
	}

	// Square ----------------------------------------------------
	g.dragSquare = newDraggableSquare()
	g.dragSquare.Build(world, g)

	g.particleSystem.SetPosition(g.dragSquare.Position().X(), g.dragSquare.Position().Y())

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.particleSystem.Update(float32(msPerUpdate))

	// Update position of particle system based on current position of rect
	// g.particleSystem.SetPosition(g.rectNode.Position().X(), g.rectNode.Position().Y())
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
	// We want the mouse events so the node can track the mouse.
	man.RegisterEventTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
	man.UnRegisterEventTarget(g)
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

func (g *gameLayer) Handle(event api.IEvent) bool {
	handled := g.dragSquare.EventHandle(event)

	g.particleSystem.SetPosition(g.dragSquare.Position().X(), g.dragSquare.Position().Y())

	if event.GetType() == api.IOTypeKeyboard {
		// fmt.Println(event.GetKeyScan())
		switch event.GetKeyScan() {
		case 65: // A
			if event.GetState() == 1 {
				g.autoTriggerEnable = !g.autoTriggerEnable
				g.particleSystem.SetAutoTrigger(g.autoTriggerEnable)
			}
		case 84: // T
			if event.GetState() == 1 {
				g.particleSystem.TriggerExplosion()
			}
		}
	}

	return handled
}
