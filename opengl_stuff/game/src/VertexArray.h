#pragma once
class VertexArray
{
public:
    VertexArray(const float* vertices, unsigned int vertexCount, const unsigned int* indices, unsigned int indexCount);
    ~VertexArray();

    void Activate();

    unsigned int GetIndexCount() const { return m_IndexCount; }
    unsigned int GetVertexCount() const { return m_VertexCount; }

private:
    unsigned int m_VertexCount;
    unsigned int m_IndexCount;
    unsigned int m_VertexBuffer;
    unsigned int m_IndexBuffer;
    unsigned int m_VertexArray;
};
