package com.gemini;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.time.Duration;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.functions.AggregateFunction;
import org.apache.flink.api.common.functions.FilterFunction;
import org.apache.flink.api.common.functions.FlatMapFunction;
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.configuration.ConfigConstants;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.connector.base.DeliveryGuarantee;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.api.windowing.assigners.TumblingProcessingTimeWindows;

public class StoreAggregationFlinkApp
{

    private static final ObjectMapper mapper = new ObjectMapper();

    public static void main(String[] args) throws Exception
    {

        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        KafkaSource<String> source = KafkaSource.<String>builder()
                                         .setBootstrapServers("kafka:29092")
                                         .setTopics("order_db.order_schema.order")
                                         .setGroupId("flink-app-group")
                                         .setStartingOffsets(OffsetsInitializer.earliest())
                                         .setValueOnlyDeserializer(new SimpleStringSchema())
                                         .build();

        KafkaSink<String> sink = KafkaSink.<String>builder()
                                     .setBootstrapServers("kafka:29092")
                                     .setRecordSerializer(KafkaRecordSerializationSchema.builder()
                                                              .setTopic("flink-out-topic")
                                                              .setValueSerializationSchema(new SimpleStringSchema())
                                                              .build())
                                     .setDeliveryGuarantee(DeliveryGuarantee.AT_LEAST_ONCE)
                                     .build();

        env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka input")
            .filter((FilterFunction<String>)(value) -> {
                try
                {
                    JsonNode rootNode = mapper.readTree(value);
                    JsonNode dataNode = mapper.readTree(rootNode.get("data").asText());
                    String status = dataNode.at("/order/order_status").asText();
                    return "Delivered".equals(status);
                }
                catch (IOException e)
                {
                    e.printStackTrace();
                    return false;
                }
            })
            .flatMap((FlatMapFunction<String, OrderData>)(value, out) -> {
                try
                {
                    JsonNode rootNode = mapper.readTree(value);
                    JsonNode dataNode = mapper.readTree(rootNode.get("data").asText());
                    JsonNode orderNode = dataNode.get("order");
                    JsonNode storeNode = dataNode.get("store");
                    if (orderNode != null && storeNode != null)
                    {
                        OrderData orderData = new OrderData(
                            orderNode.get("order_amount").asDouble(),
                            storeNode.get("store_id").asText(),
                            storeNode.at("/store_loc/store_lat").asDouble(),
                            storeNode.at("/store_loc/store_long").asDouble());
                        out.collect(orderData);
                    }
                }
                catch (IOException e)
                {
                    e.printStackTrace();
                }
            })
            // required due to java type erasure
            .returns(OrderData.class)
            .keyBy(order -> order.storeId)
            .windowAll(TumblingProcessingTimeWindows.of(Duration.ofSeconds(10)))
            .aggregate(new StoreAggregateFunction())
            .map(mapper::writeValueAsString)
            .sinkTo(sink);

        env.execute("Flink Store Aggregation Job");
    }

    public static class OrderData
    {
        @JsonProperty
        public double amount;
        @JsonProperty
        public String storeId;
        @JsonProperty
        public double lat;
        @JsonProperty
        public double lng;

        public OrderData()
        {
        }

        public OrderData(double amount, String storeId, double lat, double lng)
        {
            this.amount = amount;
            this.storeId = storeId;
            this.lat = lat;
            this.lng = lng;
        }
    }

    public static class StoreAggregatedData
    {
        @JsonProperty
        public String store_id;
        @JsonProperty
        public int order_count;
        @JsonProperty
        public double total_order_amount;
        @JsonProperty
        public double lat;
        @JsonProperty
        public double lng;

        public StoreAggregatedData()
        {
        }

        public StoreAggregatedData(String store_id, int order_count, double total_order_amount,
                                   double lat, double lng)
        {
            this.store_id = store_id;
            this.order_count = order_count;
            this.total_order_amount = total_order_amount;
            this.lat = lat;
            this.lng = lng;
        }
    }

    public static class StoreAggregateFunction
        implements AggregateFunction<OrderData, StoreAggregatedData, StoreAggregatedData>
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
