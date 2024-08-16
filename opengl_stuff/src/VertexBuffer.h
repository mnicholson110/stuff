#pragma once

class VertexBuffer
{
    private:
        unsigned int m_RendererID;
    public:
        VertexBuffer(const void* data, unsigned int size);
        ~VertexBuffer();

        void Bind() const;
        void Unbind() const;       
        
        // lock/unlock to stream data during rendering
};