package com.gemini;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Arrays;
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
import org.apache.kafka.streams.kstream.*;
import org.apache.kafka.streams.state.WindowStore;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class StoreAggregationApp
{
    private static final Logger logger = LoggerFactory.getLogger(StoreAggregationApp.class);

    public static void main(String[] args)
    {
        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, "store-aggregation-app");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "kafka:29092");
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        props.put(StreamsConfig.NUM_STREAM_THREADS_CONFIG, 3);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.String().getClass());
        // commit more frequently (or less, depending on your goal)
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000); // 1 second
        // producer overrides - send smaller batches more quickly
        // props.put(StreamsConfig.producerPrefix("linger.ms"), 10);       // default is 0 or 1 ms
        // props.put(StreamsConfig.producerPrefix("batch.size"), 32768);   // 32 KB
        // consumer overrides - process more (or fewer) records at once
        // props.put(StreamsConfig.consumerPrefix("max.poll.records"), 500);

        StreamsBuilder builder = new StreamsBuilder();
        OrderDataSerde orderDataSerde = new OrderDataSerde();
        StoreAggregatedDataSerde storeAggregatedDataSerde = new StoreAggregatedDataSerde();
        RepartitionSerde repartitionSerde = new RepartitionSerde();

        builder.stream("order_db.order_schema.order", Consumed.with(Serdes.String(), orderDataSerde))
            .filter((key, value) -> "Delivered".equals(value.order_status))
            .groupBy((key, value) -> value.store_id, Grouped.with(Serdes.String(), repartitionSerde))
            .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofSeconds(10)))
            .aggregate(StoreAggregatedData::new,
                       new StoreAggregatedDataAdder(),
                       Materialized.<String, StoreAggregatedData, WindowStore<Bytes, byte[]>>
                               as("store-aggregate-window-store")
                               .withValueSerde(storeAggregatedDataSerde))
            .toStream()
            .selectKey((windowedKey, value) -> windowedKey.key())
            .to("aggregated_store_orders", Produced.with(Serdes.String(), storeAggregatedDataSerde));

        KafkaStreams streams = new KafkaStreams(builder.build(), props);
        streams.start();
        Runtime.getRuntime().addShutdownHook(new Thread(streams::close));
    }

    public static class OrderData
    {
        public double order_amount;
        public String store_id;
        public double lat;
        public double lng;
        public String order_status;

        public OrderData()
        {
        }

        public OrderData(double order_amount, String store_id, double lat, double lng, String order_status)
        {
            this.order_amount = order_amount;
            this.store_id = store_id;
            this.lat = lat;
            this.lng = lng;
            this.order_status = order_status;
        }
    }

    public static class OrderDataSerde implements Serde<OrderData>
    {
        private static final ObjectMapper mapper = new ObjectMapper();

        public OrderDataSerde()
        {
        }

        @Override
        public Serializer<OrderData> serializer()
        {
            return new Serializer<OrderData>() {
                @Override
                public byte[] serialize(String topic, OrderData data)
                {
                    try
                    {
                        return mapper.writeValueAsBytes(data);
                    }
                    catch (IOException e)
                    {
                        logger.error("LOGGER: Serialization failed for topic {}", topic, e);
                        throw new RuntimeException("Serialization failed for topic " + topic, e);
                    }
                }
            };
        }

        @Override
        public Deserializer<OrderData> deserializer()
        {
            return new Deserializer<OrderData>() {
                @Override
                public OrderData deserialize(String key, byte[] value)
                {
                    try {
                        JsonNode rootNode = mapper.readTree(value);
                        JsonNode dataNode = mapper.readTree(rootNode.get("data").asText());
                        JsonNode orderNode = dataNode.get("order");
                        JsonNode storeNode = dataNode.get("store");
                        return new OrderData(
                                orderNode.get("order_amount").asDouble(),
                                storeNode.get("store_id").asText(),
                                storeNode.at("/store_loc/store_lat").asDouble(),
                                storeNode.at("/store_loc/store_long").asDouble(),
                                orderNode.get("order_status").asText());
                    }
                    catch (Throwable e)
                    {
                        logger.error("LOGGER: Deserialization failed for key {}", key, e);
                        return new OrderData();
                    }
                }
            };
        }
    }

    public static class RepartitionSerde implements Serde<OrderData>
    {
        private static final ObjectMapper mapper = new ObjectMapper();

        public RepartitionSerde()
        {
        }

        @Override
        public Serializer<OrderData> serializer()
        {
            return new Serializer<OrderData>() {
                @Override
                public byte[] serialize(String topic, OrderData data)
                {
                    try
                    {
                        return mapper.writeValueAsBytes(data);
                    }
                    catch (IOException e)
                    {
                        logger.error("LOGGER: Serialization failed for topic {}", topic, e);
                        throw new RuntimeException("Serialization failed for topic " + topic, e);
                    }
                }
            };
        }

        @Override
        public Deserializer<OrderData> deserializer()
        {
            return new Deserializer<OrderData>() {
                @Override
                public OrderData deserialize(String key, byte[] value)
                {
                    try {
                        JsonNode rootNode = mapper.readTree(value);
                        return new OrderData(
                                rootNode.get("order_amount").asDouble(),
                                rootNode.get("store_id").asText(),
                                rootNode.get("lat").asDouble(),
                                rootNode.get("lng").asDouble(),
                                rootNode.get("order_status").asText());
                    }
                    catch (Throwable e)
                    {
                        logger.error("LOGGER: Deserialization failed for key {}", key, e);
                        return new OrderData();
                    }
                }
            };
        }
    }

    public static class StoreAggregatedData
    {
        public String store_id;
        public int order_count;
        public double total_order_amount;
        public double lat;
        public double lng;

        public StoreAggregatedData()
        {
        }

        public StoreAggregatedData(String store_id, int order_count, double total_order_amount, double lat, double lng)
        {
            this.store_id = store_id;
            this.order_count = order_count;
            this.total_order_amount = total_order_amount;
            this.lat = lat;
            this.lng = lng;
        }
    }

    public static class StoreAggregatedDataSerde implements Serde<StoreAggregatedData>
    {
        public static final ObjectMapper mapper = new ObjectMapper();

        public StoreAggregatedDataSerde()
        {
        }

        public Serializer<StoreAggregatedData> serializer()
        {
            return new Serializer<StoreAggregatedData>() {
                @Override
                public byte[] serialize(String topic, StoreAggregatedData data)
                {
                    try
                    {
                        return mapper.writeValueAsBytes(data);
                    }
                    catch (IOException e)
                    {
                        logger.error("LOGGER: Serialization failed for topic {}", topic, e);
                        throw new RuntimeException("Serialization failed for topic " + topic, e);
                    }
                }
            };
        }

        @Override
        public Deserializer<StoreAggregatedData> deserializer()
        {
            return new Deserializer<StoreAggregatedData>() {
                @Override
                public StoreAggregatedData deserialize(String key, byte[] value)
                {
                    {
                        try
                        {
                            JsonNode rootNode = mapper.readTree(value);
                            return new StoreAggregatedData(
                                rootNode.get("store_id").asText(),
                                rootNode.get("order_count").asInt(),
                                rootNode.get("total_order_amount").asDouble(),
                                rootNode.get("lat").asDouble(),
                                rootNode.get("lng").asDouble());
                        }
                        catch (IOException e)
                        {
                            logger.error("LOGGER: Deserialization failed for key {}", key, e);
                            return new StoreAggregatedData();
                        }
                    }
                }
            };
        }
    }

    public static class StoreAggregatedDataAdder implements Aggregator<String, OrderData, StoreAggregatedData>
    {
        StoreAggregatedDataAdder()
        {
        }

        @Override
        public StoreAggregatedData apply(String s, OrderData orderData, StoreAggregatedData storeAggregatedData)
        {
            return new StoreAggregatedData(orderData.store_id,
                                           storeAggregatedData.order_count + 1,
                                           storeAggregatedData.total_order_amount + orderData.order_amount,
                                           orderData.lat,
                                           orderData.lng);
        }
    }
}
