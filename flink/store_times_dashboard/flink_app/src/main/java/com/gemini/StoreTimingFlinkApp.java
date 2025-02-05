package com.gemini;

import java.io.IOException;
import java.time.Duration;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.functions.AggregateFunction;
import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.api.common.typeinfo.TypeInformation;
import org.apache.flink.connector.base.DeliveryGuarantee;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.shaded.jackson2.com.fasterxml.jackson.databind.JsonNode;
import org.apache.flink.shaded.jackson2.com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.api.windowing.assigners.SlidingProcessingTimeWindows;

public class StoreTimingFlinkApp
{
    public static void main(String[] args) throws Exception
    {
        final String kafkaSourceAddr = System.getenv("KAFKA_SOURCE_ADDR");
        final String kafkaSinkAddr = System.getenv("KAFKA_SINK_ADDR");
        final String consumerGroupId = System.getenv("KAFKA_CONS_GROUP_ID");
        final String inputTopic = System.getenv("INPUT_TOPIC");
        final String outputTopic = System.getenv("OUTPUT_TOPIC");

        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        env.setParallelism(3);

        KafkaSource<OrderData> source = KafkaSource.<OrderData>builder()
                                            .setBootstrapServers(kafkaSourceAddr)
                                            .setTopics(inputTopic)
                                            .setGroupId(consumerGroupId)
                                            .setStartingOffsets(OffsetsInitializer.latest())
                                            .setValueOnlyDeserializer(new OrderDataDeserializationSchema())
                                            .build();

        KafkaSink<StoreAggregatedData> sink = KafkaSink.<StoreAggregatedData>builder()
                                                  .setBootstrapServers(kafkaSinkAddr)
                                                  .setRecordSerializer(KafkaRecordSerializationSchema.builder()
                                                                           .setTopic(outputTopic)
                                                                           .setKeySerializationSchema(StoreAggregatedData::serializeKey)
                                                                           .setValueSerializationSchema(StoreAggregatedData::serializeValue)
                                                                           .build())
                                                  .setDeliveryGuarantee(DeliveryGuarantee.NONE)
                                                  .build();

        env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka input")
            .filter((order) -> "Delivered".equals(order.orderStatus) || "Shipped".equals(order.orderStatus) || "Processing".equals(order.orderStatus))
            .keyBy((order) -> order.storeId)
            .window(SlidingProcessingTimeWindows.of(Duration.ofSeconds(30), Duration.ofSeconds(1)))
            .aggregate(new StoreAggregateFunction())
            .sinkTo(sink);

        env.execute("Flink Store Aggregation Job");
    }

    public static class OrderData
    {
        public String storeId;
        public String orderStatus;
        public double lat;
        public double lng;
        public Instant createdAt;
        public Instant updatedAt;

        public OrderData()
        {
        }

        public OrderData(String storeId, String orderStatus, double lat, double lng, Instant createdAt, Instant updatedAt)
        {
            this.storeId = storeId;
            this.orderStatus = orderStatus;
            this.lat = lat;
            this.lng = lng;
            this.createdAt = createdAt;
            this.updatedAt = updatedAt;
        }
    }

    public static class OrderDataDeserializationSchema implements DeserializationSchema<OrderData>
    {
        public static final ObjectMapper mapper = new ObjectMapper();

        @Override
        public OrderData deserialize(byte[] message)
        {
            try
            {
                JsonNode rootNode = mapper.readTree(message);
                JsonNode dataNode = mapper.readTree(rootNode.get("data").asText());
                JsonNode orderNode = dataNode.get("order");
                JsonNode storeNode = dataNode.get("store");
                return new OrderData(
                    storeNode.get("store_id").asText(),
                    orderNode.get("order_status").asText(),
                    storeNode.at("/store_loc/store_lat").asDouble(),
                    storeNode.at("/store_loc/store_long").asDouble(),
                    Instant.parse(rootNode.get("created_at").asText()),
                    Instant.parse(rootNode.get("updated_at").asText()));
            }
            catch (IOException e)
            {
                e.printStackTrace();
                return new OrderData();
            }
        }

        @Override
        public boolean isEndOfStream(OrderData o)
        {
            return false;
        }

        @Override
        public TypeInformation<OrderData> getProducedType()
        {
            return TypeInformation.of(OrderData.class);
        }
    }

    public static class StoreAggregatedData
    {
        public String store_id;
        public int count;
        public double avg_processing;
        public double avg_shipped;
        public double avg_delivered;
        public double lat;
        public double lng;

        public static final ObjectMapper mapper = new ObjectMapper();

        public StoreAggregatedData()
        {
        }

        public byte[] serializeValue()
        {
            try
            {
                return mapper.writeValueAsString(this).getBytes();
            }
            catch (Throwable e)
            {
                e.printStackTrace();
                return new byte[0];
            }
        }

        public byte[] serializeKey()
        {
            return store_id.getBytes();
        }

        public StoreAggregatedData(String store_id, int count, double avg_processing, double avg_shipped, double avg_delivered, double lat, double lng)
        {
            this.store_id = store_id;
            this.count = count;
            this.avg_processing = avg_processing;
            this.avg_shipped = avg_shipped;
            this.avg_delivered = avg_delivered;
            this.lat = lat;
            this.lng = lng;
        }
    }

    public static class StoreAggregateFunction implements AggregateFunction<OrderData, StoreOrderStatusAccumulator, StoreAggregatedData>
    {
        public StoreAggregateFunction()
        {
        }

        @Override
        public StoreOrderStatusAccumulator createAccumulator()
        {
            return new StoreOrderStatusAccumulator();
        }

        @Override
        public StoreOrderStatusAccumulator add(OrderData orderData, StoreOrderStatusAccumulator accumulator)
        {
            long durationMillis = ChronoUnit.MILLIS.between(orderData.createdAt, orderData.updatedAt);

            if (accumulator.storeId == null)
            {
                accumulator.storeId = orderData.storeId;
                accumulator.lat = orderData.lat;
                accumulator.lng = orderData.lng;
            }

            if ("Processing".equals(orderData.orderStatus))
            {
                accumulator.processingSum += durationMillis;
                accumulator.processingCount++;
            }
            else if ("Shipped".equals(orderData.orderStatus))
            {
                accumulator.shippedSum += durationMillis;
                accumulator.shippedCount++;
            }
            else if ("Delivered".equals(orderData.orderStatus))
            {
                accumulator.deliveredSum += durationMillis;
                accumulator.deliveredCount++;
            }

            return accumulator;
        }

        @Override
        public StoreAggregatedData getResult(StoreOrderStatusAccumulator accumulator)
        {
            double avgProcessing = (accumulator.processingCount > 0)
                                       ? ((double)accumulator.processingSum / accumulator.processingCount)
                                       : 0.0;
            double avgShipped = (accumulator.shippedCount > 0)
                                    ? ((double)accumulator.shippedSum / accumulator.shippedCount)
                                    : 0.0;
            double avgDelivered = (accumulator.deliveredCount > 0)
                                      ? ((double)accumulator.deliveredSum / accumulator.deliveredCount)
                                      : 0.0;

            return new StoreAggregatedData(accumulator.storeId, accumulator.processingCount, avgProcessing, avgShipped, avgDelivered, accumulator.lat, accumulator.lng);
        }

        @Override
        public StoreOrderStatusAccumulator merge(StoreOrderStatusAccumulator a, StoreOrderStatusAccumulator b)
        {
            StoreOrderStatusAccumulator merged = new StoreOrderStatusAccumulator();
            merged.processingCount = a.processingCount + b.processingCount;
            merged.processingSum = a.processingSum + b.processingSum;
            merged.shippedCount = a.shippedCount + b.shippedCount;
            merged.shippedSum = a.shippedSum + b.shippedSum;
            merged.deliveredCount = a.deliveredCount + b.deliveredCount;
            merged.deliveredSum = a.deliveredSum + b.deliveredSum;
            return merged;
        }
    }

    public static class StoreOrderStatusAccumulator
    {
        public String storeId = null;
        public double lat;
        public double lng;
        public int processingCount = 0;
        public int processingSum = 0;
        public int shippedCount = 0;
        public int shippedSum = 0;
        public int deliveredCount = 0;
        public int deliveredSum = 0;
    }
}
