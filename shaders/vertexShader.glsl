#version 410 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColour;
layout (location = 2) in vec2 aTexCoord;

out vec3 ourColour;
out vec2 TexCoord;

uniform float scale;
uniform mat4 translation;
uniform mat4 projection;

void main()
{
    gl_Position = projection * translation * vec4(aPos, 1.0);
    ourColour = aColour;
    TexCoord = aTexCoord;
}
