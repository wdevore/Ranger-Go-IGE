#version 450 core

// This input attribute handles the stream of vertices sent to this shader
// The client defines the structure layout using vertex-attribute-pointer
layout (location = 0) in vec3 aPos;

uniform mat4 model;

// These uniforms don't change and are set once at the start of the client App
uniform mat4 view;
uniform mat4 projection;

void main()
{
    // note that we read the multiplication from right to left
    gl_Position = projection * view * model * vec4(aPos, 1.0);
}
