#version 410 core

in vec3 ourColour;

out vec4 FragColor;

void main()
{
    FragColor = vec4(ourColour, 1.0);
} 
