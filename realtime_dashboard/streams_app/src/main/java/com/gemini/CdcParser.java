package com.gemini;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

public class CdcParser {
    private static final ObjectMapper mapper = new ObjectMapper();

    public static ParsedOrder parseCdcEvent(CdcEvent event) {
        if (event == null || event.getData() == null) return null;

        try {
            // "data" is a JSON string. Let's parse it into a JsonNode:
            JsonNode root = mapper.readTree(event.getData());

            // Extract order fields
            JsonNode orderNode = root.path("order");
            double amount = orderNode.path("order_amount").asDouble(0.0);
            String status = orderNode.path("order_status").asText("");

            // Extract store fields
            JsonNode storeNode = root.path("store");
            long storeId = storeNode.path("store_id").asLong(0);
            JsonNode locNode = storeNode.path("store_loc");
            double lat = locNode.path("store_lat").asDouble(0.0);
            double lng = locNode.path("store_long").asDouble(0.0);
            String addr = storeNode.path("store_addr").asText("");

            // Build the ParsedOrder
            ParsedOrder parsed = new ParsedOrder();
            parsed.setStoreId(storeId);
            parsed.setStoreLat(lat);
            parsed.setStoreLong(lng);
            parsed.setStoreAddr(addr);
            parsed.setOrderAmount(amount);
            parsed.setOrderStatus(status);

            return parsed;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }
}
