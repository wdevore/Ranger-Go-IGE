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
Scenes are stack based.
Scene are stored in a pool. Any scene that elects not to be placed back in the pool is destroyed.

The Boot scene is special. It is the only scene that appears immediately without transitioning onto the stage. As the boot scene runs a status is sent back to the node manager (NM) indicating if it is still busy.
When the NM receives a "done" signal from the boot scene it immediately pops the stack for the next outgoing scene.

At this point there are two scenes to run: the outgoing just popped and the incoming scene on the stack top.

The incoming scene can't begin entering the stage until the outgoing scene has completed its task. Once complete it sends signal to node manager.

The outgoing scene can either remain on stage until the incoming scene has completed its transition or it can transition off in sync with the incoming scene. At this point the outgoing scene either vanishes or is pooled. For example, boot scene and splash scene both vanish.

If there is no incoming scene then when the outgoing scene exits the game exits.

#### **Scenes**
Each scene waits to enter the stage. The NM monitors the current scene on stage and notifies the waiting scene.

Scenes waiting to enter the stage are queued on a stack. Once the last scene has pulled from the stack and finishes the game is *over*.

Scenes are only pulled from the stack when the current scene signals *SceneTransitioningOut* or *SceneFinished*.

If current scene signaled *SceneTransitioningOut* then the next scene on the stack is brought in as the **incoming** scene and the *EnterNode()* is called on the incoming scene. The incoming scene then begins to transition onto the stage. Once it has completed transitioning it signals NM which then moves it to currentScene.

If a scene needs to direct to another scene, for example *Menu* ```-->``` *Settings*.
* *Menu* first pushes itself onto the **stack**.
* NM then pushes *Settings* onto stack.
* *Menu* then tells NM that is it *SceneTransitioningOut* and that a ***return*** scene is queued after the stack-top.
* This causes NM to make *Settings* the incoming scene via *EnterNode()*
* When *Settings* indicates that it is *SceneTransitioningOut* NM pops it into the running scene and activates the stack-top (aka Menu)

It is the current scene that determines what scene becomes active next. For example, *Boot* specifies that *Splash* is next by placing *Splash* on the stack, and *Splash* specifies *Menu* is next.

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

While the boot scene runs it sends its status to NM. Status are *SceneBusy* and *SceneFinished*.

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
* --On each Push NM calls the scene's *EnterNode()* for any prep work.
* The NM pops *Boot* and makes it the current-scene and issues *EnterStageNode()* at which point Boot can start entering the stage.
* *Boot* begins issuing *SceneTransitioningIn* signal while it transitions.
* Once transition is complete *Boot* begins issuing *SceneBusy*.
* When *Boot* is done it issues *SceneFinished* rather than *SceneTransitioningOut*.
* NM detects *SceneFinished* and tells the incoming (Splash) scene that it can start transitioning onto the stage by calling *EnterStageNode()*.
* *Splash* now begins issuing *SceneTransitioningIn* as it begins transitioning onto the stage.
* *Splash* is now running as the current-scene and begins issuing *SceneBusy*.
* When *Splash* is finished and is ready to leave the stage it issues *SceneFinished*.
* NM detects this and begins running the incoming-scene (Menu) by sending *EnterStageNode()* to *Menu*.
* Both *Splash* and *Menu* are now running.
* *Menu* begins issuing *SceneTransitioningIn* and *Splash* begins issuing *SceneTransitioningOut*.
* Once the incoming-scene (aka Menu) finihes transitioning issues *SceneTransitionComplete*.
* NM makes current-scene = incoming-scene which effectively **discards** the previous current-scene (Splash) even if it hasn't completed any transitioning. No point in running a scene that can't be seen.
* *Menu* begins issuing *SceneBusy* signal while it is running.
* User selects **Settings**.
* *Menu* pushes itself onto the stack.
* *Menu* pushes *Settings* onto the stack which causes *Settings*'s *EnterNode()* to be called. *Settings* needs to track state if it has been entered prior.
* *Menu* begins issuing *SceneTransitioningOut* and *Settings* begins issuing *SceneTransitioningIn*.
* Once *Settings* finishes transitioning onto the stage it issues *SceneTransitionComplete*.
* NM detect and makes current-scene = incoming-scene (Settings) and *Menu* is discarded but is held on the stack.
* User selects to exit *Settings* scene. *Settings* scene issues the *SceneTransitioningOut*.
* NM activates stack-top and sends *EnterStageNode()* to *Menu*.
* Both *Settings* and *Menu* are now running.
* *Menu* finishes transitioning and issues *SceneTransitionComplete*.
* *Menu* is now running.
---


* Once Boot completes its task it notifies NM, NM then notifies Splash. At the same time NM pops Splash and continues to run it, NM also begins running what is at the top which is Menu.
* Once Splash is complete it notifies NM and may begin transitioning off the stage.
* NM then notifies the scene on the top of the stack and pops it.
* Menu comes off the stack as the outgoing scene and transitions onto the stage.
* At the same time 
* The Menu scene will typically push itself back on the stack after popping the next scene to run.