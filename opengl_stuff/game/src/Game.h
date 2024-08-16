#pragma once
#include <SDL2/SDL.h>
#include <GL/glew.h>
#include <cstdint>

class Game
{
public:
    Game();

    bool Init();
    void RunLoop();
    void Shutdown();

    void CreateVertexArray();
    void CreateShader();

private:
    void ProcessInput();
    void Update();
    void Render();

    SDL_Window* m_Window;
    SDL_GLContext m_Context;
    bool m_IsRunning;
    Uint32 m_TicksCount;
    class VertexArray* m_VertexArray;
    class Shader* m_Shader;
};