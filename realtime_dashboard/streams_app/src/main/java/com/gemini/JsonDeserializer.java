package com.gemini;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.kafka.common.serialization.Deserializer;

/**
 * Generic JSON Deserializer using Jackson.
 */
public class JsonDeserializer<T> implements Deserializer<T> {

    private final ObjectMapper objectMapper = new ObjectMapper();
    private Class<T> targetType;

    public JsonDeserializer(Class<T> targetType) {
        this.targetType = targetType;
    }

    @Override
    public T deserialize(String topic, byte[] data) {
        if (data == null) {
            return null;
        }
        try {
            return objectMapper.readValue(data, targetType);
        } catch (Exception e) {
            throw new RuntimeException("Error deserializing JSON message", e);
        }
    }
}
