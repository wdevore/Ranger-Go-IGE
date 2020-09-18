# Ranger-Go-IGE
**Ranger-Go-IGE** (IGE) is a continuation of [RangerGO](https://github.com/wdevore/RangerGo).

As was RangerGo, so is IGE, a variation of the [Ranger Dart](https://github.com/wdevore/Ranger-Dart) game engine but written in [Go](https://golang.org/) but using [OpenGL](https://www.opengl.org/) and [GLFW](https://www.glfw.org/).

This version is the great merging of RangerGo, RangerGo-GLFW and Ranger-Alpha.

# Current Tasks and Goals
* **Done** OpenGL 4.x Core
* **Done** Pixel Font (Slow version)
* **Done** Pixel Font (Faster version)
* **Done** Static shapes: Circle
* **Done** Dynamic shapes: Line
* **Done** Space mappings
* **Done** add Outline shapes
* **Done** Alpha coloring / transparency
* **Done** Node Dragging
* **Done** Zoom Node
* **Done** Filters: transform and translate
* **Done** Particles
* **Done** Animation (tweening) *Using tanema's library*: https://github.com/tanema/gween
* **Done** Box2D physics (with space ship)
  * Uses: https://github.com/ByteArena/box2d
* **Done** Zones combined with Zoom
* **Done** Interpolation
* **Done** Sprite Textures (quads)
* **Done** Scene Transitioning
* **Done** Bitmap fonts
* **Done** Simple motion animations
* **Done** Audio (SFXR 8bit sound: https://sfxr.me/) (However, 16bit stills needs completion)
  * built GUI using imGui: https://github.com/inkyblackness/imgui-go
  * build using: https://github.com/faiface/beep
* <b style="color:red">*working*</b> QuadTree for view-space culling
  * Example 1: Random spread of points/particles into tree (with capacity)
  * Example 2: Random spread of squares
  * Example 3: Insert squares at Mouse
  * Example 4: Triangle ship moving through tree
  * Example 5: Query and highlight tree according to obj-obj intersecting
* Simple Widget GUI framework
  * Buttons
  * Checkboxes
  * ListBox and Combo dropdowns
  * Inputs (text, float, int, bool)
  * Text
  * Dialogs (OK, Yes/No, etc.)
  * Grouping (i.e. Radio buttons)
* Physics with Textures
* Joysticks and Gamepads
* OpenGL ES (https://github.com/golang/go/wiki/Mobile)
* Shaders with interleved vertex/color, for example, checkboards
* Stippling with OpenGL patterns (aka advanced shaders)
* Custom vector font (needs lower case completion)
* Vector Fonts via SVG import
