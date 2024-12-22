package com.gemini;

public class StoreAggregation {

    private long storeId;
    private long orderCount;
    private double totalOrderAmount;
    private double storeLat;
    private double storeLong;
    private String storeAddr;

    // No-arg constructor (required by some frameworks like Jackson)
    public StoreAggregation() {
    }

    // All-arg constructor
    public StoreAggregation(long storeId, long orderCount, double totalOrderAmount,
                            double storeLat, double storeLong, String storeAddr) {
        this.storeId = storeId;
        this.orderCount = orderCount;
        this.totalOrderAmount = totalOrderAmount;
        this.storeLat = storeLat;
        this.storeLong = storeLong;
        this.storeAddr = storeAddr;
    }

    // Getters and Setters

    public long getStoreId() {
        return storeId;
    }

    public void setStoreId(long storeId) {
        this.storeId = storeId;
    }

    public long getOrderCount() {
        return orderCount;
    }

    public void setOrderCount(long orderCount) {
        this.orderCount = orderCount;
    }

    public double getTotalOrderAmount() {
        return totalOrderAmount;
    }

    public void setTotalOrderAmount(double totalOrderAmount) {
        this.totalOrderAmount = totalOrderAmount;
    }

    public double getStoreLat() {
        return storeLat;
    }

    public void setStoreLat(double storeLat) {
        this.storeLat = storeLat;
    }

    public double getStoreLong() {
        return storeLong;
    }

    public void setStoreLong(double storeLong) {
        this.storeLong = storeLong;
    }

    public String getStoreAddr() {
        return storeAddr;
    }

    public void setStoreAddr(String storeAddr) {
        this.storeAddr = storeAddr;
    }

    @Override
    public String toString() {
        return "StoreAggregation{" +
                "storeId=" + storeId +
                ", orderCount=" + orderCount +
                ", totalOrderAmount=" + totalOrderAmount +
                ", storeLat=" + storeLat +
                ", storeLong=" + storeLong +
                ", storeAddr='" + storeAddr + '\'' +
                '}';
    }
}
