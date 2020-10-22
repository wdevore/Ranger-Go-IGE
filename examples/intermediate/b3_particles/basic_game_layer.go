package main

import (
	"fmt"
	"math/rand"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/particles"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
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
	// Instead of using two node: vline and hline, I'm using one "+ node.
	xyAxis, err := shapes.NewMonoPlusNode("XYAxis", world, world.Underlay())
	if err != nil {
		return err
	}
	xyAxis.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	ghl := xyAxis.(*shapes.MonoPlusNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	tri, err := shapes.NewMonoTriangleNode("Tri", api.FILLOUTLINED, world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100.0)
	tri.SetPosition(150.0, 0.0)
	gsq := tri.(*shapes.MonoTriangleNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.PanSkin))

	// Particle system
	activator := particles.NewActivator360()
	g.particleSystem = particles.NewParticleSystem(activator)
	g.particleSystem.Activate(true)
	g.particleSystem.SetAutoTrigger(false)

	colors := []api.IPalette{
		color.NewPaletteInt64(0xFBD872FF),
		color.NewPaletteInt64(0xFFC845FF),
		color.NewPaletteInt64(0xFFB81CFF),
		color.NewPaletteInt64(0xC69214FF),
		color.NewPaletteInt64(0xAD841FFF),
	}

	// Now populate the system
	for i := 0; i < 50; i++ {
		ptri, err := shapes.NewMonoCircleNode(fmt.Sprintf("%s%d", "::Part", i), api.FILLED, 5, world, g)
		if err != nil {
			return err
		}
		gptri := ptri.(*shapes.MonoCircleNode)
		ci := int(rand.Float32() * 4)
		gptri.SetFilledColor(colors[ci])
		gptri.SetVisible(false)
		gptri.SetScale(15.0)
		p := particles.NewNodeParticle(ptri)
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

	// Update position of particle system based on current position of square
	g.particleSystem.SetPosition(g.dragSquare.Position().X(), g.dragSquare.Position().Y())
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
