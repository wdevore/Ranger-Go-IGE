#version 450 core

out vec4 color;

uniform vec4 fragColor;

void main()
{
    color = vec4(fragColor);
}
