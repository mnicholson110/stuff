#pragma once
#include <GL/glew.h>
#include <string>
#include <glm/glm.hpp>

class Shader
{
public:
    Shader();
    ~Shader();

    void Activate();
    void SetMatrixUniform(const char* name, const glm::mat4& matrix);

private:
    bool CompileShader(const char* fileName, GLenum shaderType, GLuint& outShader);
    bool IsCompiled(GLuint shader);
    bool IsValidProgram();

    GLuint m_VertexShader;
    GLuint m_FragmentShader;
    GLuint m_Program;

public:
    const char* vertexShader =
        "#version 450 core\n"
        "layout(location = 0) in vec4 position;\n"
        "void main()\n"
        "{\n"
        "   gl_Position = position;\n"
        "}\n";

    const char* fragmentShader =
        "#version 450 core\n"
        "layout(location = 0) out vec4 color;\n"
        "void main()\n"
        "{\n"
        "   color = vec4(1.0, 0.0, 0.0, 1.0);\n"
        "}\n";
};