#version 410 core

in vec4 ourColour;
in vec2 TexCoord;

uniform sampler2D ourTexture;

out vec4 FragColor;

void main()
{
    FragColor = mix(texture(ourTexture, TexCoord), ourColour, 0.0);
} 
