# Ranger-Go-IGE

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
https://learnopengl.com/Getting-started/Hello-Triangle
https://www.reddit.com/r/opengl/comments/3515bi/rendering_multiple_objects_from_multiple_vaos/

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
http://quabr.com:8182/41789384/go-gl-rendering-vbo-not-displaying
https://github.com/YagoCarballo/Go-GL-Assignment-2

## Easing
https://www.shadertoy.com/view/Xd2yRd Alternative easing functions

## Knowledge
https://learnopengl.com/Getting-started/Coordinate-Systems
http://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
https://www.haroldserrano.com/blog/loading-vertex-normal-and-uv-data-onto-opengl-buffers

## SVG

https://github.com/JoshVarga/svgparser
https://play.golang.org/p/kyfff6Kg1c
https://github.com/rustyoz/svg
https://golang.hotexamples.com/examples/github.com.catiepg.svg/-/Parse/golang-parse-function-examples.html

## Bitmap fonts
https://www.gamedev.net/tutorials/_/technical/opengl/opengl-texture-mapping-an-introduction-r947/
https://www.opengl.org/archives/resources/code/samples/glut_examples/texfont/texfont.html
http://plib.sourceforge.net/fnt/index.html
https://fontforge.org/en-US/
https://learn.adafruit.com/custom-fonts-for-pyportal-circuitpython-display/conversion
https://lazyfoo.net/tutorials/OpenGL/20_bitmap_fonts/index.php

http://nadev.zapto.org/2019/05/27/creating-a-bitmap-font/ part #1
http://nadev.zapto.org/2019/05/27/split-linear-font-file/ part #2
http://nadev.zapto.org/2019/03/29/combining-letter-images/ part #3

https://en.wikipedia.org/wiki/Glyph_Bitmap_Distribution_Format .bdf format

https://gimplearn.net/viewtopic.php?f=4&t=317&p=1513 gimp script

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
https://github.com/McNopper
https://github.com/McNopper/OpenGL
https://user.xmission.com/~nate/glut.html
https://www.desultoryquest.com/blog/drawing-anti-aliased-circular-points-using-opengl-slash-webgl/ anti-aliased points

### Shaders
https://www.geeks3d.com/hacklab/20190225/demo-checkerboard-in-glsl/
https://stackoverflow.com/questions/4694608/glsl-checkerboard-pattern
https://thebookofshaders.com/09/

## Mobile (Android)
https://github.com/golang/go/wiki/Mobile