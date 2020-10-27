# Ranger-Go-IGE
**Ranger-Go-IGE** (IGE) is a continuation of [RangerGO](https://github.com/wdevore/RangerGo).

As was RangerGo, so is IGE, a variation of the [Ranger Dart](https://github.com/wdevore/Ranger-Dart) game engine but written in [Go](https://golang.org/) but using [OpenGL](https://www.opengl.org/) and [GLFW](https://www.glfw.org/).

This version is the great merging of RangerGo, RangerGo-GLFW and Ranger-Alpha.

# Games
* Port Spout
* Port StarCastle / Star Barriers
* Port Bosconian
* Port Omega Race
* Port 2048
* Create New Game. A blend of **Lunar Lander** (Sky and Landings), **Moon Patrol** (Ground), **Atari Pitfall** (Obstacles), **Defender** (Rescue). Multi-planet and Moons.

# Demos
* Combine vector fonts with tweening and physics. The text rotates and then randomly falls to the ground. Then the text moves back inline using tweening.

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
* <b style="color:orange">Partial Completion</b> QuadTrees
  * Example 1: Random spread of points/particles into tree (with capacity)
  * Example 2: Random spread of squares
  * Example 3: Insert squares at Mouse
  * Example 4: Triangle ship moving through tree
  * **Done** Example 5: Query and highlight tree according to obj-obj intersecting
* <b style="color:red">*working*</b> Custom vector font (needs lower case completion)
  * Vector Fonts via SVG import
* Improve Bitmap fonts. Rename current bitmap font object to "CharacterBitmapFont" and add new one called "StringBitmapFont". The new one creates a buffer of quads.
* Scripting via https://wren.io/ or make custom language using via http://craftinginterpreters.com/
* Simple Widget GUI framework
  * Buttons
  * Checkboxes
  * ListBox and Combo dropdowns
  * Inputs (text, float, int, bool)
  * Text
  * Dialogs (OK, Yes/No, etc.)
  * Grouping (i.e. Radio buttons)
* Batch rendering
* OpenGL ES (https://github.com/golang/go/wiki/Mobile)
* Physics with Textures
* Joysticks and Gamepads
* Shaders with interleved vertex/color, for example, checkboards
* Stippling with OpenGL patterns (aka advanced shaders)
