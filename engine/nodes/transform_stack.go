package nodes

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

type transformStackItem struct {
	current api.IAffineTransform
}

func newTransformItem() *transformStackItem {
	o := new(transformStackItem)
	o.current = maths.NewTransform()
	return o
}

type transformStack struct {
	stack    []*transformStackItem
	stackTop int

	current api.IAffineTransform
	post    api.IAffineTransform // Pre allocated cache
}

const transformStackDepth = 100

func newTransformStack() *transformStack {
	o := new(transformStack)

	o.current = maths.NewTransform()
	o.post = maths.NewTransform()

	return o
}

func (t *transformStack) Initialize() {
	t.stack = make([]*transformStackItem, transformStackDepth)

	for i := 0; i < transformStackDepth; i++ {
		t.stack[i] = newTransformItem()
	}

	// Apply centered view-space matrix
	// rc.Apply(rc.world.ViewSpace())
}

func (t *transformStack) Apply(aft api.IAffineTransform) {
	// Concat this transform onto the current transform but don't push it.
	// Use post multiply
	maths.Multiply(aft, t.current, t.post)
	t.current.SetByTransform(t.post)
}

func (t *transformStack) Save() {
	top := t.stack[t.stackTop]
	top.current.SetByTransform(t.current)

	t.stackTop++
}

func (t *transformStack) Restore() {
	t.stackTop--

	top := t.stack[t.stackTop]

	t.current.SetByTransform(top.current)
}
