# Ranger-Go-IGE
**Ranger-Go-IGE** (IGE) is a continuation of [RangerGO](https://github.com/wdevore/RangerGo).

As was RangerGo, so is IGE, a variation of the [Ranger Dart](https://github.com/wdevore/Ranger-Dart) game engine but written in [Go](https://golang.org/) but using [OpenGL](https://www.opengl.org/) and [GLFW](https://www.glfw.org/).

This version is the great merging of RangerGo, RangerGo-GLFW and Ranger-Alpha.

# Current Tasks and Goals
* **Working** OpenGL 4.x and perhaps (ES)
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
The are two config files: one for the engine itself and the other for each game.
 
## Articles
* [tutorial-opengl-with-golang](https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl)
* [box2d iforce2d](https://www.iforce2d.net/b2dapps/)
* https://github.com/go-gl/glfw

# Tracking (Optional)
Some **Nodes**/**Objects** may want to *Track* the properties of another **Node**.

For example, an AABB object may wan't to track **Mesh** changes on a node such that it can *rebuild* its internal min/max properties.

## Packages

```
go get github.com/tanema/gween
go get github.com/ByteArena/box2d
go get github.com/go-gl/gl/v4.5-core/gl
go get github.com/go-gl/glfw/v3.3/glfw
go get -u github.com/go-gl/gl/v4.5-{core,compatibility}/gl
```

# VS Code setup:
"window.zoomLevel": 1.9
"editor.fontSize": 14
"editor.fontLigatures": true
"terminal.integrated.fontSize": 12
"editor.fontFamily": "'Cascadia Code', 'Cascadia Mono', 'Courier New', monospace, Consolas"

# Example and source
https://github.com/go-gl/mathgl GLM
https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
https://github.com/cstegel/opengl-samples-golang
https://github.com/crockeo/go-tuner
https://github.com/pikkpoiss/ld33
https://github.com/pikkpoiss/twodee
https://github.com/esenti/godraw
https://github.com/runningwild/glop
https://github.com/btmura/blockcillin keys  
https://github.com/manueldun/Game
https://github.com/kurrik/opengl-benchmarks
https://github.com/tbogdala/cubez physics engine
https://github.com/Triangle345/GT
https://github.com/yucchiy/toybox-opengl
https://github.com/CubeLite/gltext-1
https://github.com/cozely/cozely
https://github.com/LonelyPale/go-opengl
