#version 450 core

in vec2 TexCoords;
out vec4 color;

uniform sampler2DArray textureArry;
uniform int layer;

void main()
{
    color = texture(textureArry, vec3(TexCoords.xy, layer));
}
