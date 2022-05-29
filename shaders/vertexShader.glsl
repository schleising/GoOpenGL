#version 410 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec4 aColour;
layout (location = 2) in vec2 aTexCoord;

out vec4 ourColour;
out vec2 TexCoord;

uniform mat4 scale;
uniform mat4 rotation;
uniform mat4 translation;
uniform mat4 projection;

void main()
{
    gl_Position = projection * translation * rotation * scale * vec4(aPos, 1.0);
    ourColour = aColour;
    TexCoord = aTexCoord;
}
