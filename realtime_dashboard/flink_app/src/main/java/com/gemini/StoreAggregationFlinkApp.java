package com.gemini;

import java.io.IOException;
import java.time.Duration;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.functions.AggregateFunction;
import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.api.common.serialization.SerializationSchema;
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

public class StoreAggregationFlinkApp
{
    public static void main(String[] args) throws Exception
    {
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        env.setParallelism(3);

        KafkaSource<OrderData> source = KafkaSource.<OrderData>builder()
                                            .setBootstrapServers("kafka:29092")
                                            .setTopics("order_db.order_schema.order")
                                            .setGroupId("flink-app-group")
                                            .setStartingOffsets(OffsetsInitializer.latest())
                                            .setValueOnlyDeserializer(new OrderDataDeserializationSchema())
                                            .build();

        KafkaSink<StoreAggregatedData> sink = KafkaSink.<StoreAggregatedData>builder()
                                                  .setBootstrapServers("kafka:29092")
                                                  .setRecordSerializer(KafkaRecordSerializationSchema.builder()
                                                                           .setTopic("aggregated_store_orders_flink")
                                                                           .setKeySerializationSchema(StoreAggregatedData::serializeKey)
                                                                           .setValueSerializationSchema(new StoreAggregatedDataSerializationSchema())
                                                                           .build())
                                                  .setDeliveryGuarantee(DeliveryGuarantee.NONE)
                                                  .build();

        env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka input")
            .filter((order) -> "Delivered".equals(order.orderStatus))
            .keyBy((order) -> order.storeId)
            .window(SlidingProcessingTimeWindows.of(Duration.ofSeconds(10), Duration.ofSeconds(1)))
            .aggregate(new StoreAggregateFunction())
            .sinkTo(sink);

        env.execute("Flink Store Aggregation Job");
    }

    public static class OrderData
    {
        public double amount;
        public String storeId;
        public double lat;
        public double lng;
        public String orderStatus;

        public OrderData()
        {
        }

        public OrderData(double amount, String storeId, double lat, double lng, String orderStatus)
        {
            this.amount = amount;
            this.storeId = storeId;
            this.lat = lat;
            this.lng = lng;
            this.orderStatus = orderStatus;
        }
    }

    public static class OrderDataDeserializationSchema implements DeserializationSchema<OrderData>
    {
        public final ObjectMapper mapper = new ObjectMapper();
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
                    orderNode.get("order_amount").asDouble(),
                    storeNode.get("store_id").asText(),
                    storeNode.at("/store_loc/store_lat").asDouble(),
                    storeNode.at("/store_loc/store_long").asDouble(),
                    orderNode.get("order_status").asText());
            }
            catch (IOException e)
            {
                e.printStackTrace();
                return new OrderData(0.0, null, 0.0, 0.0, null);
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
        public int order_count;
        public double total_order_amount;
        public double lat;
        public double lng;

        public final ObjectMapper mapper = new ObjectMapper();

        public StoreAggregatedData()
        {
        }

        public byte[] serializeKey()
        {
            return store_id.getBytes();
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

    public static class StoreAggregatedDataSerializationSchema implements SerializationSchema<StoreAggregatedData>
    {
        final ObjectMapper mapper = new ObjectMapper();

        @Override
        public byte[] serialize(StoreAggregatedData data) {
            try
            {
                return mapper.writeValueAsString(data).getBytes();
            }
            catch (IOException e)
            {
                e.printStackTrace();
                return new byte[0];
            }
        }
    }

    public static class StoreAggregateFunction implements AggregateFunction<OrderData, StoreAggregatedData, StoreAggregatedData>
    {
        @Override
        public StoreAggregatedData createAccumulator()
        {
            return new StoreAggregatedData("", 0, 0.0, 0.0, 0.0);
        }

        @Override
        public StoreAggregatedData add(OrderData orderData, StoreAggregatedData accumulator)
        {
            if (accumulator.store_id == null || accumulator.store_id.isEmpty())
            {
                accumulator.store_id = orderData.storeId;
                accumulator.lat = orderData.lat;
                accumulator.lng = orderData.lng;
            }
            accumulator.order_count += 1;
            accumulator.total_order_amount += orderData.amount;
            return accumulator;
        }

        @Override
        public StoreAggregatedData getResult(StoreAggregatedData accumulator)
        {
            return accumulator;
        }

        @Override
        public StoreAggregatedData merge(StoreAggregatedData a, StoreAggregatedData b)
        {
            a.order_count += b.order_count;
            a.total_order_amount += b.total_order_amount;
            return a;
        }
    }
}
