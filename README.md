# Ranger-Go-IGE
**Ranger-Go-IGE** (IGE) is a continuation of [RangerGO](https://github.com/wdevore/RangerGo).

As was RangerGo so is IGE, a variation of the [Ranger Dart](https://github.com/wdevore/Ranger-Dart) game engine but written in [Go](https://golang.org/) but using [OpenGL](https://www.opengl.org/).

# Current Tasks and Goals
* **Done** Node Dragging
* **Done** Filters: transform and translate
* **Done** Zoom Node
* **Done** Interpolation
* **Done** Simple motion animations
* **Done** Circle, AABB
* **Done** AnchorNode
* **Done** Particles
* **Done** Animation (tweening) -- Using tanema's library
* **Done** Box2D physics (with space ship)
* **Done** Zones combined with Zoom
* **Working** OpenGL 4.x and perhaps (ES)
* Audio (SFXR 8bit sound)
  * May build GUI using imGui
  * May use: https://github.com/faiface/beep
* QuadTree for view-space culling
* Simple Widget GUI framework
  * Buttons
  * Checkboxes
  * Text input
  * Text
  * Dialog
  * Grouping (i.e. Radio buttons)
* Sprite Textures ()
* Finish lower case Vector font characters.
* Enhance raster fonts to allow transforms
* Joysticks and Gamepads

## Notes

# Tracking (Optional)
Some **Nodes**/**Objects** may want to *Track* the properties of another **Node**.

For example, an AABB object may wan't to track **Mesh** changes on a node such that it can *rebuild* its internal min/max properties.

## Packages

```
go get github.com/tanema/gween
go get github.com/ByteArena/box2d
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.2/glfw
go get -u github.com/go-gl/gl/v4.1-{core,compatibility}/gl
```
