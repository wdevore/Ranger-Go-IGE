#version 450 core

in vec2 TexCoords;
out vec4 color;

uniform sampler2D texture1;
uniform vec4 fragColor;

void main()
{
    color = texture(texture1, TexCoords) * vec4(fragColor);
    // color = texture(texture1, TexCoords) * vec4(1.0, 0.5, 0.0, 1.0);
    // color = vec4(1.0, 0.5, 0.0, 1.0);
}
