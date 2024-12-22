package com.gemini;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.*;
import org.apache.kafka.streams.kstream.*;
import org.apache.kafka.streams.state.KeyValueStore;

import com.gemini.CdcEvent;
import com.gemini.ParsedOrder;
import com.gemini.StoreAggregation;
import com.gemini.JsonSerializer;
import com.gemini.JsonDeserializer;

import java.util.Properties;

/**
 * Main Kafka Streams application that filters for Delivered orders,
 * aggregates by store_id, and emits store-level stats for real-time use.
 */
public class StoreAggregationApp {


    // ------------------------------------------------------------------------
    // 5) The MAIN method: Building the Streams topology
    // ------------------------------------------------------------------------
    public static void main(String[] args) {

        // Kafka Streams configuration
        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, "store-aggregation-app");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "kafka:29092");
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        StreamsBuilder builder = new StreamsBuilder();

        // 5.1) Read from the Debezium CDC topic
        String inputTopic = "order_db.order_schema.order";
        String outputTopic = "aggregated_store_orders";

        // KStream of raw JSON messages (key=String, value=String)
        KStream<String, String> rawStream = builder.stream(
            inputTopic,
            Consumed.with(Serdes.String(), Serdes.String())
        );

        // 5.2) Parse top-level JSON -> CdcEvent
        KStream<String, CdcEvent> cdcEvents = rawStream.mapValues(value -> parseCdcEvent(value));

        // 5.3) Convert CdcEvent -> ParsedOrder (extract the nested "data" field)
        KStream<String, ParsedOrder> parsedOrders = cdcEvents.mapValues(event -> parseNestedData(event));

        // 5.4) Filter for only "Delivered" orders
        KStream<String, ParsedOrder> deliveredOrders = parsedOrders
            .filter((key, parsed) -> parsed != null
                    && "Delivered".equalsIgnoreCase(parsed.getOrderStatus()));

        // 5.5) Re-key by storeId, then group
        KGroupedStream<Long, ParsedOrder> groupedByStore = deliveredOrders
            .selectKey((key, parsed) -> parsed.getStoreId())  // new key = storeId
            .groupByKey(Grouped.with(Serdes.Long(), getParsedOrderSerde()));

        // 5.6) Aggregate into a StoreAggregation
        KTable<Long, StoreAggregation> storeAggTable = groupedByStore.aggregate(
            // Initializer
            () -> new StoreAggregation(0, 0, 0.0, 0.0, 0.0, ""),

            // Aggregation logic
            (storeId, newOrder, aggValue) -> {
                aggValue.setStoreId(storeId);
                aggValue.setOrderCount(aggValue.getOrderCount() + 1);
                aggValue.setTotalOrderAmount(aggValue.getTotalOrderAmount() + newOrder.getOrderAmount());
                // Update store info (assuming 1 storeId always has same location & address)
                aggValue.setStoreLat(newOrder.getStoreLat());
                aggValue.setStoreLong(newOrder.getStoreLong());
                aggValue.setStoreAddr(newOrder.getStoreAddr());
                return aggValue;
            },

            // Materialize
            Materialized.<Long, StoreAggregation, KeyValueStore<org.apache.kafka.common.utils.Bytes, byte[]>>as("store-agg-store")
                .withKeySerde(Serdes.Long())
                .withValueSerde(getStoreAggregationSerde())
        );

        // 5.7) Produce the aggregated results to outputTopic
        storeAggTable.toStream()
                     .to(outputTopic, Produced.with(Serdes.Long(), getStoreAggregationSerde()));

        // Start Kafka Streams
        KafkaStreams streams = new KafkaStreams(builder.build(), props);
        streams.start();

        // Graceful shutdown
        Runtime.getRuntime().addShutdownHook(new Thread(streams::close));
    }

    // ------------------------------------------------------------------------
    // 6) Utility methods for parsing & Serdes
    // ------------------------------------------------------------------------

    // Parse the top-level CDC message into CdcEvent
    private static CdcEvent parseCdcEvent(String json) {
        if (json == null) return null;
        try {
            return new ObjectMapper().readValue(json, CdcEvent.class);
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    // Parse the nested "data" JSON from CdcEvent into ParsedOrder
    private static ParsedOrder parseNestedData(CdcEvent event) {
        if (event == null || event.getData() == null) return null;
        try {
            ObjectMapper mapper = new ObjectMapper();
            JsonNode root = mapper.readTree(event.getData());
            JsonNode orderNode = root.path("order");
            JsonNode storeNode = root.path("store");

            // Extract order fields
            double amount = orderNode.path("order_amount").asDouble(0.0);
            String status = orderNode.path("order_status").asText("");

            // Extract store fields
            long storeId = storeNode.path("store_id").asLong(0);
            String addr = storeNode.path("store_addr").asText("");
            JsonNode loc = storeNode.path("store_loc");
            double lat = loc.path("store_lat").asDouble(0.0);
            double lng = loc.path("store_long").asDouble(0.0);

            // Build ParsedOrder
            ParsedOrder parsed = new ParsedOrder();
            parsed.setOrderAmount(amount);
            parsed.setOrderStatus(status);
            parsed.setStoreId(storeId);
            parsed.setStoreAddr(addr);
            parsed.setStoreLat(lat);
            parsed.setStoreLong(lng);

            return parsed;

        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    // Serde for ParsedOrder
    private static Serde<ParsedOrder> getParsedOrderSerde() {
        return Serdes.serdeFrom(new JsonSerializer<>(), new JsonDeserializer<>(ParsedOrder.class));
    }

    // Serde for StoreAggregation
    private static Serde<StoreAggregation> getStoreAggregationSerde() {
        return Serdes.serdeFrom(new JsonSerializer<>(), new JsonDeserializer<>(StoreAggregation.class));
    }
}
