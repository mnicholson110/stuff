package com.gemini;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.time.Duration;
import java.util.Collections;
import java.util.List;
import java.util.Properties;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.Deserializer;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.common.serialization.Serializer;
import org.apache.kafka.common.utils.Bytes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.Grouped;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.streams.kstream.Produced;
import org.apache.kafka.streams.kstream.TimeWindows;
import org.apache.kafka.streams.state.WindowStore;

public class StoreAggregationApp
{
    private static final ObjectMapper mapper = new ObjectMapper();

    public static void main(String[] args)
    {
        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, "store-aggregation-app");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "kafka:29092");
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.String().getClass());
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass());
        props.put(StreamsConfig.NUM_STREAM_THREADS_CONFIG, 3);
        // commit more frequently (or less, depending on your goal)
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 500); // 1 second
        // producer overrides - send smaller batches more quickly
        // props.put(StreamsConfig.producerPrefix("linger.ms"), 10);       // default is 0 or 1 ms
        // props.put(StreamsConfig.producerPrefix("batch.size"), 32768);   // 32 KB
        // consumer overrides - process more (or fewer) records at once
        // props.put(StreamsConfig.consumerPrefix("max.poll.records"), 500);

        StreamsBuilder builder = new StreamsBuilder();
        JsonSerde<StoreAggregatedData> serde = new JsonSerde<>(StoreAggregatedData.class);

        builder.stream("order_db.order_schema.order", Consumed.with(Serdes.String(), Serdes.String()))
            .filter((key, value) -> {
                try
                {
                    JsonNode rootNode = mapper.readTree(value);
                    JsonNode dataNode = mapper.readTree(rootNode.get("data").asText());
                    return "Delivered".equals(dataNode.at("/order/order_status").asText());
                }
                catch (IOException e)
                {
                    e.printStackTrace();
                    return false;
                }
            })
            .flatMapValues(value -> {
                try
                {
                    JsonNode rootNode = mapper.readTree(value);
                    JsonNode dataNode = mapper.readTree(rootNode.get("data").asText());
                    JsonNode orderNode = dataNode.get("order");
                    JsonNode storeNode = dataNode.get("store");
                    return List.of(mapper.writeValueAsString(new OrderData(
                        orderNode.get("order_amount").asDouble(),
                        storeNode.get("store_id").asText(),
                        storeNode.at("/store_loc/store_lat").asDouble(),
                        storeNode.at("/store_loc/store_long").asDouble())));
                }
                catch (IOException e)
                {
                    e.printStackTrace();
                    return Collections.emptyList();
                }
            })
            .groupBy((key, orderData) -> {
                try
                {
                    OrderData data = mapper.readValue(orderData, OrderData.class);
                    return data.storeId;
                }
                catch (IOException e)
                {
                    e.printStackTrace();
                    return null;
                }
            }, Grouped.with(Serdes.String(), Serdes.String()))
            .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofSeconds(10)))
            .aggregate(StoreAggregatedData::new, (key, value, aggregate) -> {
                try
                {
                    OrderData orderData = mapper.readValue(value, OrderData.class);
                    return new StoreAggregatedData(
                        orderData.storeId,
                        aggregate.order_count + 1,
                        aggregate.total_order_amount + orderData.amount,
                        orderData.lat,
                        orderData.lng);
                }
                catch (IOException e)
                {
                    e.printStackTrace();
                    return aggregate;
                }
            }, Materialized.<String, StoreAggregatedData, WindowStore<Bytes, byte[]>>as("store-aggregate-window-store").withKeySerde(Serdes.String()).withValueSerde(serde))
            .toStream()
            .map((windowedKey, aggregatedValue) -> KeyValue.pair(windowedKey.key(), aggregatedValue))
            .to("aggregated_store_orders", Produced.with(Serdes.String(), serde));

        KafkaStreams streams = new KafkaStreams(builder.build(), props);
        streams.start();
        Runtime.getRuntime().addShutdownHook(new Thread(streams::close));
    }

    static class OrderData
    {
        @JsonProperty
        double amount;
        @JsonProperty
        String storeId;
        @JsonProperty
        double lat;
        @JsonProperty
        double lng;

        OrderData()
        {
        }

        OrderData(double amount, String storeId, double lat, double lng)
        {
            this.amount = amount;
            this.storeId = storeId;
            this.lat = lat;
            this.lng = lng;
        }
    }

    static class StoreAggregatedData
    {
        @JsonProperty
        String store_id;
        @JsonProperty
        int order_count;
        @JsonProperty
        double total_order_amount;
        @JsonProperty
        double lat;
        @JsonProperty
        double lng;

        StoreAggregatedData()
        {
        }

        StoreAggregatedData(String store_id, int order_count, double total_order_amount, double lat, double lng)
        {
            this.store_id = store_id;
            this.order_count = order_count;
            this.total_order_amount = total_order_amount;
            this.lat = lat;
            this.lng = lng;
        }
    }

    static class JsonSerde<T> implements Serde<T>
    {
        private final Class<T> clazz;
        private final ObjectMapper mapper = new ObjectMapper();

        public JsonSerde(Class<T> clazz)
        {
            this.clazz = clazz;
        }

        @Override
        public Serializer<T> serializer()
        {
            return new Serializer<T>() {
                @Override
                public byte[] serialize(String topic, T data)
                {
                    try
                    {
                        return mapper.writeValueAsBytes(data);
                    }
                    catch (IOException e)
                    {
                        throw new RuntimeException("Serialization failed for topic " + topic, e);
                    }
                }
            };
        }

        @Override
        public Deserializer<T> deserializer()
        {
            return new Deserializer<T>() {
                @Override
                public T deserialize(String topic, byte[] data)
                {
                    if (data == null)
                    {
                        return null;
                    }
                    try
                    {
                        return mapper.readValue(data, clazz);
                    }
                    catch (IOException e)
                    {
                        throw new RuntimeException("Deserialization failed for topic " + topic, e);
                    }
                }
            };
        }
    }
}
