#include "Actor.h"

Actor::Actor(Game* game)
    : m_Game(game),
    m_State(stateActive),
    m_Position(0.0f, 0.0f),
    m_Scale(1.0f),
    m_Rotation(0.0f)
{
}

Actor::~Actor() {}

void Actor::Update(float deltaTime) {}

void Actor::ProcessInput(const uint8_t* keyState)
{
}
