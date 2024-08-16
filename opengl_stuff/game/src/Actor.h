#pragma once

#include <glm/glm.hpp>

class Actor
{
public:
    enum State
    {
        stateActive,
        statePaused,
        stateDead
    };

    Actor(class Game* game);
    virtual ~Actor();

    virtual void Update(float deltaTime);

    virtual void ProcessInput(const uint8_t* keyState);

    const State& GetState() const { return m_State; }
    void SetState(State state) { m_State = state; }

    const glm::vec2& GetPosition() const { return m_Position; }
    void SetPosition(const glm::vec2& position) { m_Position = position; }

    const float& GetScale() const { return m_Scale; }
    void SetScale(float scale) { m_Scale = scale; }

    const float& GetRotation() const { return m_Rotation; }
    void SetRotation(float rotation) { m_Rotation = rotation; }

    void ComputeTransform();
    const glm::mat4& GetTransform() const { return m_Transform; }



    class Game* GetGame() const { return m_Game; }

private:
    State m_State;

    glm::vec2 m_Position;
    float m_Scale;
    float m_Rotation;
    glm::mat4 m_Transform;
    bool m_recomputeTransform;

    class Game* m_Game;
};
