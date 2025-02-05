package com.gemini;

import java.io.IOException;
import java.time.Duration;
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

public class PaymentsFlinkApp
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

        KafkaSource<PaymentData> source = KafkaSource.<PaymentData>builder()
                                              .setBootstrapServers(kafkaSourceAddr)
                                              .setTopics(inputTopic)
                                              .setGroupId(consumerGroupId)
                                              .setStartingOffsets(OffsetsInitializer.latest())
                                              .setValueOnlyDeserializer(new PaymentDataDeserializationSchema())
                                              .build();

        KafkaSink<CustomerAggregatedData> sink = KafkaSink.<CustomerAggregatedData>builder()
                                                     .setBootstrapServers(kafkaSinkAddr)
                                                     .setRecordSerializer(KafkaRecordSerializationSchema.builder()
                                                                              .setTopic(outputTopic)
                                                                              .setKeySerializationSchema(CustomerAggregatedData::serializeKey)
                                                                              .setValueSerializationSchema(CustomerAggregatedData::serializeValue)
                                                                              .build())
                                                     .setDeliveryGuarantee(DeliveryGuarantee.NONE)
                                                     .build();

        env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka input")
            .keyBy((payment) -> payment.customer_name)
            .window(SlidingProcessingTimeWindows.of(Duration.ofSeconds(20), Duration.ofSeconds(1)))
            .aggregate(new CustomerAggregateFunction())
            .filter((agg) -> agg.payment_count > 3)
            .sinkTo(sink);

        env.execute("Flink Store Aggregation Job");
    }

    public static class PaymentData
    {
        public int auth_id;
        public int store_id;
        public String customer_name;
        public int last_four;

        public PaymentData()
        {
        }

        public PaymentData(int auth_id, int store_id, String customer_name, int last_four)
        {
            this.auth_id = auth_id;
            this.store_id = store_id;
            this.customer_name = customer_name;
            this.last_four = last_four;
        }
    }

    public static class PaymentDataDeserializationSchema implements DeserializationSchema<PaymentData>
    {
        public static final ObjectMapper mapper = new ObjectMapper();

        @Override
        public PaymentData deserialize(byte[] message)
        {
            try
            {
                JsonNode rootNode = mapper.readTree(message);
                return new PaymentData(
                    rootNode.get("auth_id").asInt(),
                    rootNode.get("store_id").asInt(),
                    rootNode.get("customer_name").asText(),
                    rootNode.get("last_four").asInt());
            }
            catch (IOException e)
            {
                e.printStackTrace();
                return new PaymentData();
            }
        }

        @Override
        public boolean isEndOfStream(PaymentData o)
        {
            return false;
        }

        @Override
        public TypeInformation<PaymentData> getProducedType()
        {
            return TypeInformation.of(PaymentData.class);
        }
    }

    public static class CustomerAggregatedData
    {
        public String customer_name;
        public int payment_count;

        public static final ObjectMapper mapper = new ObjectMapper();

        public CustomerAggregatedData()
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
            return customer_name.getBytes();
        }

        public CustomerAggregatedData(String customer_name, int payment_count)
        {
            this.customer_name = customer_name;
            this.payment_count = payment_count;
        }
    }

    public static class CustomerAggregateFunction implements AggregateFunction<PaymentData, CustomerAggregatedData, CustomerAggregatedData>
    {
        public CustomerAggregateFunction()
        {
        }

        @Override
        public CustomerAggregatedData createAccumulator()
        {
            return new CustomerAggregatedData();
        }

        @Override
        public CustomerAggregatedData add(PaymentData paymentData, CustomerAggregatedData accumulator)
        {
            if (accumulator.customer_name == null)
            {
                accumulator.customer_name = paymentData.customer_name;
            }
            accumulator.payment_count += 1;
            return accumulator;
        }

        @Override
        public CustomerAggregatedData getResult(CustomerAggregatedData accumulator)
        {
            return accumulator;
        }

        @Override
        public CustomerAggregatedData merge(CustomerAggregatedData a, CustomerAggregatedData b)
        {
            a.payment_count += b.payment_count;
            return a;
        }
    }
}
