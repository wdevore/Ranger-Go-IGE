#version 450 core

out vec4 color;

uniform vec3 fragColor;

void main()
{
    color = vec4(fragColor, 1.0f);
}