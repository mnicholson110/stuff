package com.gemini;

public class ParsedOrder {
    private long storeId;
    private double storeLat;
    private double storeLong;
    private String storeAddr;
    private double orderAmount;
    private String orderStatus;

    public ParsedOrder() {}

    // Getters and setters

    public long getStoreId() {
        return storeId;
    }

    public void setStoreId(long storeId) {
        this.storeId = storeId;
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

    public double getOrderAmount() {
        return orderAmount;
    }

    public void setOrderAmount(double orderAmount) {
        this.orderAmount = orderAmount;
    }

    public String getOrderStatus() {
        return orderStatus;
    }

    public void setOrderStatus(String orderStatus) {
        this.orderStatus = orderStatus;
    }
}
