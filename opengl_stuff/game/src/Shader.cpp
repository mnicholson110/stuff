#include "Shader.h"
#include <SDL2/SDL.h>

Shader::Shader()
    : m_VertexShader(0),
    m_FragmentShader(0),
    m_Program(0)
{
    if (!CompileShader(vertexShader, GL_VERTEX_SHADER, m_VertexShader)
        || !CompileShader(fragmentShader, GL_FRAGMENT_SHADER, m_FragmentShader))
    {
        SDL_Log("Failed to compile shader");
        return;
    }

    m_Program = glCreateProgram();
    glAttachShader(m_Program, m_VertexShader);
    glAttachShader(m_Program, m_FragmentShader);
    glLinkProgram(m_Program);

    if (!IsValidProgram())
    {
        SDL_Log("Failed to link program");
        return;
    }
}

Shader::~Shader()
{
    glDeleteProgram(m_Program);
    glDeleteShader(m_VertexShader);
    glDeleteShader(m_FragmentShader);
}

void Shader::Activate()
{
    glUseProgram(m_Program);
}

void Shader::SetMatrixUniform(const char* name, const glm::mat4& matrix)
{
}

bool Shader::CompileShader(const char* shaderSource, GLenum shaderType, GLuint& outShader)
{
    outShader = glCreateShader(shaderType);
    glShaderSource(outShader, 1, &shaderSource, nullptr);
    glCompileShader(outShader);

    if (!IsCompiled(outShader))
    {
        SDL_Log("Failed to compile shader");
        return false;
    }
    return true;
}

bool Shader::IsCompiled(GLuint shader)
{
    GLint status;
    glGetShaderiv(shader, GL_COMPILE_STATUS, &status);

    if (status != GL_TRUE)
    {
        char buffer[512];
        memset(buffer, 0, 512);
        glGetShaderInfoLog(shader, 511, nullptr, buffer);
        SDL_Log("GLSL Compile Failed:\n%s", buffer);
        return false;
    }

    return true;
}

bool Shader::IsValidProgram()
{

    GLint status;
    glGetProgramiv(m_Program, GL_LINK_STATUS, &status);
    if (status != GL_TRUE)
    {
        char buffer[512];
        memset(buffer, 0, 512);
        glGetProgramInfoLog(m_Program, 511, nullptr, buffer);
        SDL_Log("GLSL Link Status:\n%s", buffer);
        return false;
    }

    return true;
}
