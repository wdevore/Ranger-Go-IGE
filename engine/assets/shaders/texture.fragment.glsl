#version 450 core

in vec2 TexCoords;
out vec4 color;

uniform sampler2D texture1;
// uniform vec4 fragColor;

void main()
{
    // Note: fragColor must be passed to a vec4 even though it
    // is passed to the shader as a vec4. Weird.
    // vec4(fragColor) * 
    color = texture(texture1, TexCoords);
} // Note: 450-core requires a blank line at EOF
