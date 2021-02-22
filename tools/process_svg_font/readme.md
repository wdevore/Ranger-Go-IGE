Note: Inkscape's coord space is:

This tool reads an ".svg" file and parses it to produce a json font file ".vfon"

## JSON

```
{
  "A": {
    "OffsetX": 1.0,
    "OffsetY": 1.0,
    "Bounds", [0.0, 0.0, 0.0, 1.0],
    "Path": {
      "Closed": true,
      "Vertices", [1.0, 0.5, 0.0, 1.0]
    },
    "Path": {
      "Closed": true,
      "Vertices", [1.0, 0.5, 0.0, 1.0]
    }
  },
  "a": {
    "OffsetX": 1.0,
    "OffsetY": 1.0,
    "Bounds", [0.0, 0.0, 0.0, 1.0],
    "Path": {
      "Closed": true,
      "Vertices", [1.0, 0.5, 0.0, 1.0]
    },
    "Path": {
      "Closed": true,
      "Vertices", [1.0, 0.5, 0.0, 1.0]
    }
  }
}
```

In order to port the svg coord data into Ranger's coord-space we need to perform a few transforms in Inkscape-space first. For that we need:
* AABB
* Calc an Offset from text's position relative to AABB's lower-left

Eventually the final output is normalized coords in Ranger-space, including the offsets:

```
            ^ +Y
        0,0 |                     1,1
            .----------------------.
            |                      |
            |                      |
            |                      |
            |                      |
            |                      |
            |                      |
            |                      |
            |                      |
            |                      |
            .----------------------.---> +X
          0,0                     1,0
            ^
            | Offset dY
            |
            |
.---------->.----------------------|
  Offset dX        aabb.width

```

The text can either ignore the aabb.width and use fix char offsets, or use proportional by using the width + gap.

----
First we find the vertex aabb bounds.

Use the lower-left (inkscape-space) and text position to get offset delta.

Each char is shifted to the left bounds-edge and a width recorded.
Each char is converted from model to unit-space based on the largest
bounds found in the atlas.

# Fonts
* Cascadia Code
* Ubuntu Mono
* URW Gothic
* Zil Semi Slab



```
Inkscape's coord space is:

    .----------------------.----------------------.
    |                      |                      |
    |                      |                      |
    |                      |                      |
    |                      |                      |
    |       Quadrant       |       Quadrant       |
    |          1           |          2           |
    |                      |                      |
    |                      |                      |
    |                      |                      |
    .----------------------.----------------------.---> +X
    |                      | TL  |                |
    |                      |     |                |
    |                      |-----. offset         |
    |                      |                      |
    |       Quadrant       |       Quadrant       |
    |          4           |          3           |
    |                      |                      |
    |                      |                      |
    |                      |                      |
    .----------------------.----------------------.BR
                           |
                           v +Y

Ranger's coord space is:

                       ^ +Y
                       |
.----------------------.----------------------.
|                      |                      |
|                      |                      |
|                      |                      |
|                      |                      |
|       Quadrant       |       Quadrant       |
|          1           |          2           |
|                      |                      |
|                      |                      |
|                      |                      |
.----------------------.----------------------.---> +X
|                      | TL                   |
|                      |                      |
|                      |                      |
|                      |                      |
|       Quadrant       |       Quadrant       |
|          4           |          3           |
|                      |                      |
|                      |                      |
|                      |                      |
.----------------------.----------------------. BR
```

Find top-left most position and transform all vertices such that
the char moves to quadrant #3 is inkscape-space.
Then reflect about the +X axis to put the coords in quadrant #2.