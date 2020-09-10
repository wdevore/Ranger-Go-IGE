# Ranger-Go-IGE

## Working task
When switching between rendergraphics we need to consider both the shader and Vao. Some nodes only need switch one or the other or maybe both.
Each RG should know what to "use" prior to rendering. It should also know if it is already in use and not re-use redundantly.

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
* "window.zoomLevel": 1.9
* "editor.fontSize": 14
* "editor.fontLigatures": true
* "terminal.integrated.fontSize": 12
* "editor.fontFamily": "'Cascadia Code', 'Cascadia Mono', 'Courier New', monospace, Consolas"

# Example and source
* https://learnopengl.com/Getting-started/Hello-Triangle
* https://www.reddit.com/r/opengl/comments/3515bi/rendering_multiple_objects_from_multiple_vaos/
* https://github.com/go-gl/mathgl GLM
* https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
* https://github.com/cstegel/opengl-samples-golang
* https://github.com/crockeo/go-tuner
* https://github.com/pikkpoiss/ld33
* https://github.com/pikkpoiss/twodee
* https://github.com/esenti/godraw
* https://github.com/runningwild/glop
* https://github.com/btmura/blockcillin keys  
* https://github.com/manueldun/Game
* https://github.com/kurrik/opengl-benchmarks
* https://github.com/tbogdala/cubez physics engine
* https://github.com/Triangle345/GT
* https://github.com/yucchiy/toybox-opengl
* https://github.com/CubeLite/gltext-1
* https://github.com/cozely/cozely
* https://github.com/LonelyPale/go-opengl
* http://quabr.com:8182/41789384/go-gl-rendering-vbo-not-displaying
* https://github.com/YagoCarballo/Go-GL-Assignment-2

## Easing
https://www.shadertoy.com/view/Xd2yRd Alternative easing functions

## Knowledge
* https://learnopengl.com/Getting-started/Coordinate-Systems
* http://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
* https://www.haroldserrano.com/blog/loading-vertex-normal-and-uv-data-onto-opengl-buffers

## SVG

* https://github.com/JoshVarga/svgparser
* https://play.golang.org/p/kyfff6Kg1c
* https://github.com/rustyoz/svg
* https://golang.hotexamples.com/examples/github.com.catiepg.svg/-/Parse/golang-parse-function-examples.html

## Bitmap fonts
* https://en.wikibooks.org/wiki/OpenGL_Programming/Modern_OpenGL_Tutorial_Text_Rendering_02
* https://github.com/pbnjay/pixfont
* https://www.cnx-software.com/2020/06/19/fontedit-font-editor-targets-embedded-systems-with-led-lcd-or-e-paper-displays/
* https://www.gamedev.net/tutorials/_/technical/opengl/opengl-texture-mapping-an-introduction-r947/
* https://www.opengl.org/archives/resources/code/samples/glut_examples/texfont/texfont.html
* http://plib.sourceforge.net/fnt/index.html
* https://fontforge.org/en-US/
* https://learn.adafruit.com/custom-fonts-for-pyportal-circuitpython-display/conversion
* https://lazyfoo.net/tutorials/OpenGL/20_bitmap_fonts/index.php
* http://nadev.zapto.org/2019/05/27/creating-a-bitmap-font/ part #1
* http://nadev.zapto.org/2019/05/27/split-linear-font-file/ part #2
* http://nadev.zapto.org/2019/03/29/combining-letter-images/ part #3
* https://en.wikipedia.org/wiki/Glyph_Bitmap_Distribution_Format .bdf format
* https://gimplearn.net/viewtopic.php?f=4&t=317&p=1513 gimp script

### FNT
https://ttf2fnt.com/ convert ttf to fnt

FNT is a bitmap format while TTF is an outline/vector format. To get FNT output, you'll first need to create a bitmap strike or strikes. Fontforge can do this if it has been built with Freetype included:

Go to Element -> Bitmap Strikes Available.
Select the 'Win' button
Under 'Point Sizes' enter the size (or sizes) you want.
Tick 'Use FreeType' and 'Create Rasterized Strikes' Click OK.
You should now have a bitmap strike to work with. Go to Generate Fonts; on the left side select 'No Outline Font', and on the right side select Windows FNT, pick the size you want, and generate.

https://libgdx.badlogicgames.com/tools.html


## OpenGL examples and source
* https://learnopengl.com/code_viewer_gh.php?code=src/1.getting_started/4.1.textures/textures.cpp
* https://github.com/faiface/glhf
* https://github.com/McNopper
* https://github.com/McNopper/OpenGL
* https://user.xmission.com/~nate/glut.html
* https://www.desultoryquest.com/blog/drawing-anti-aliased-circular-points-using-opengl-slash-webgl/ anti-aliased points

### 2D texture arrays
* https://ferransole.wordpress.com/2014/06/09/array-textures/

### Stippling Dashes
* https://stackoverflow.com/questions/52928678/dashed-line-in-opengl3 stippling
* http://jcgt.org/published/0002/02/08/paper.pdf
* https://stackoverflow.com/questions/43392333/fragment-shader-for-stipple-pattern
* https://community.khronos.org/t/how-to-draw-a-line-stipple-with-es-2-0/72531/7
Computer graphics through opengl 4.3 pg 782 2nd ed.
* https://html.developreference.com/article/14406767/fragment+shader+for+stipple+pattern

### Shaders
* https://www.geeks3d.com/hacklab/20190225/demo-checkerboard-in-glsl/
* https://stackoverflow.com/questions/4694608/glsl-checkerboard-pattern
* https://thebookofshaders.com/09/

## Mobile (Android)
* https://github.com/golang/go/wiki/Mobile

## Scenes and transitions
Each scene waits to enter the stage. The NM monitors the current scene on stage and notifies the waiting scene.

Scenes waiting to enter the stage are queued on a stack. Once the last scene has pulled from the stack and finishes the game is *over*.

Scenes are only pulled from the stack when the current scene signals *SceneTransitioningOut* or *SceneFinished*.

If current scene signaled *SceneTransitioningOut* then the next scene on the stack is brought in as the **incoming** scene and the *EnterNode()* is called on the incoming scene. The incoming scene then begins to transition onto the stage. Once it has completed transitioning it signals NM which then moves it to currentScene.

*Menu*s can have multiple **targets** or destinations, for example, *Settings*, *HighScore*, *Game*, etc. When a user selects a destination *Menu* first pushes itself onto the stack then pushes the destination. At this point the *Menu* is both the current scene (i.e. transitioning off the stage) and is on the stack underneath the incoming scene. Both the current-scene and incoming scene are active. Once that *Menu* has finished transitioning it signals the NM and the NM pops the stack (incoming scene) into the current scene.
```
Running scene: Boot

Stack:
   Splash
   Settings <-- SubScene
   Highscore <-- SubScene
   Game
   Menu
```

#### **Boot scene**
To start a game a Boot scene must be present and appears immediately.

#### **Splash scene**
The splash scene waits off screen until the NM signals that the current scene (aka boot) is either transitioning off the stage or is finished.

#### *Use case*
Lets cover a typical case:
```
Boot --> Splash --> Menu --> Settings -->  Menu
                                             |
                                             v
          Menu <-- Game <-- Menu <-- HighScore
```


* The game pushes Exit, Menu, Splash then Boot onto the stack.
   ```
   Stack:
      Boot
      Splash
      Menu
      Exit
   ```
# Audio
Audio has both internal and external values. 
The audio generator works on values that aren't exactly meaningfull to users and that is what the external values are for.
There is a converter to translate internal to external prior to viewing by a user.
It is the internal values that are persisted.
