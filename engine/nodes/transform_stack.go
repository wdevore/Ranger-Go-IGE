package nodes

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

type transformStackItem struct {
	current api.IMatrix4
}

func newTransformItem() *transformStackItem {
	o := new(transformStackItem)
	o.current = maths.NewMatrix4()
	return o
}

type transformStack struct {
	stack    []*transformStackItem
	stackTop int

	current api.IMatrix4
	post    api.IMatrix4 // Pre allocated cache

	m4 api.IMatrix4
}

const transformStackDepth = 100

func newTransformStack() api.ITransformStack {
	o := new(transformStack)

	o.current = maths.NewMatrix4()
	o.post = maths.NewMatrix4()
	o.m4 = maths.NewMatrix4()

	return o
}

func (t *transformStack) Initialize(mat api.IMatrix4) {
	t.stack = make([]*transformStackItem, transformStackDepth)

	for i := 0; i < transformStackDepth; i++ {
		t.stack[i] = newTransformItem()
	}

	// The initial value ready for the top of the stack.
	t.current.Set(mat)
}

func (t *transformStack) Apply(aft api.IMatrix4) api.IMatrix4 {
	// Concat this transform onto the current transform but don't push it.
	// Use post multiply
	maths.Multiply4(t.current, aft, t.post) // Post-multiply
	t.current.Set(t.post)
	return t.current
}

func (t *transformStack) ApplyAffine(aft api.IAffineTransform) api.IMatrix4 {
	// Concat this transform onto the current transform but don't push it.
	// Use post multiply
	t.m4.SetFromAffine(aft)
	maths.Multiply4(t.current, t.m4, t.post) // Post-multiply
	t.current.Set(t.post)
	return t.current
}

func (t *transformStack) Save() {
	top := t.stack[t.stackTop]
	top.current.Set(t.current)
	t.stackTop++
}

func (t *transformStack) Restore() {
	t.stackTop--
	top := t.stack[t.stackTop]
	t.current.Set(top.current)
}
