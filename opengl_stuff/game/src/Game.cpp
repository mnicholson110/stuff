#include "Game.h"
#include "VertexArray.h"
#include "Shader.h"

Game::Game()
    : m_Window(nullptr),
    m_Context(nullptr),
    m_IsRunning(true),
    m_TicksCount(0),
    m_VertexArray(nullptr),
    m_Shader(nullptr)
{
}

bool Game::Init()
{
    if (SDL_Init(SDL_INIT_VIDEO) != 0)
    {
        SDL_Log("Unable to initialize SDL: %s", SDL_GetError());
        return false;
    }

    m_Window = SDL_CreateWindow(
        "OpenGL",
        SDL_WINDOWPOS_CENTERED,
        SDL_WINDOWPOS_CENTERED,
        800,
        600,
        SDL_WINDOW_OPENGL);

    if (!m_Window)
    {
        SDL_Log("Failed to create window: %s", SDL_GetError());
        return false;
    }

    m_Context = SDL_GL_CreateContext(m_Window);

    glewInit();

    CreateVertexArray();
    CreateShader();

    return true;
}

void Game::CreateVertexArray()
{
    float positions[8] = {
        -0.5f, -0.5f,
         0.5f, -0.5f,
         0.5f,  0.5f,
        -0.5f,  0.5f
    };

    unsigned int indices[6] = {
        0, 1, 2,
        2, 3, 0
    };

    m_VertexArray = new VertexArray(positions, 4, indices, 6);
}

void Game::CreateShader()
{
    m_Shader = new Shader();
}

void Game::RunLoop()
{
    while (m_IsRunning)
    {
        ProcessInput();
        Update();
        Render();
    }
}

void Game::Shutdown()
{
    delete m_VertexArray;
    delete m_Shader;
    SDL_GL_DeleteContext(m_Context);
    SDL_DestroyWindow(m_Window);
    SDL_Quit();
}

void Game::ProcessInput()
{
    SDL_Event event;
    while (SDL_PollEvent(&event))
    {
        switch (event.type)
        {
        case SDL_QUIT:
            m_IsRunning = false;
            break;
        }
    }

    const Uint8* keyState = SDL_GetKeyboardState(NULL);
    if (keyState[SDL_SCANCODE_ESCAPE])
    {
        m_IsRunning = false;
    }
}

void Game::Update()
{
}

void Game::Render()
{
    m_Shader->Activate();

    glClearColor(0.0f, 0.0f, 0.0f, 1.0f);
    glClear(GL_COLOR_BUFFER_BIT);

    glDrawElements(GL_TRIANGLES, 6, GL_UNSIGNED_INT, nullptr);

    SDL_GL_SwapWindow(m_Window);
}