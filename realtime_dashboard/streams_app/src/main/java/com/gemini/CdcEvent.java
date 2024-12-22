package com.gemini;

public class CdcEvent {

    private long order_id;
    private String data;
    private String created_at;
    private String updated_at;
    private String __op;
    private String __table;
    private long __lsn;
    private long __source_ts_ms;

    // No-argument constructor (required by some frameworks like Jackson)
    public CdcEvent() {
    }

    // All-argument constructor
    public CdcEvent(long order_id, String data, String created_at, String updated_at,
                    String __op, String __table, long __lsn, long __source_ts_ms) {
        this.order_id = order_id;
        this.data = data;
        this.created_at = created_at;
        this.updated_at = updated_at;
        this.__op = __op;
        this.__table = __table;
        this.__lsn = __lsn;
        this.__source_ts_ms = __source_ts_ms;
    }

    // Getters and Setters

    public long getOrder_id() {
        return order_id;
    }

    public void setOrder_id(long order_id) {
        this.order_id = order_id;
    }

    public String getData() {
        return data;
    }

    public void setData(String data) {
        this.data = data;
    }

    public String getCreated_at() {
        return created_at;
    }

    public void setCreated_at(String created_at) {
        this.created_at = created_at;
    }

    public String getUpdated_at() {
        return updated_at;
    }

    public void setUpdated_at(String updated_at) {
        this.updated_at = updated_at;
    }

    public String get__op() {
        return __op;
    }

    public void set__op(String __op) {
        this.__op = __op;
    }

    public String get__table() {
        return __table;
    }

    public void set__table(String __table) {
        this.__table = __table;
    }

    public long get__lsn() {
        return __lsn;
    }

    public void set__lsn(long __lsn) {
        this.__lsn = __lsn;
    }

    public long get__source_ts_ms() {
        return __source_ts_ms;
    }

    public void set__source_ts_ms(long __source_ts_ms) {
        this.__source_ts_ms = __source_ts_ms;
    }

    @Override
    public String toString() {
        return "CdcEvent{" +
                "order_id=" + order_id +
                ", data='" + data + '\'' +
                ", created_at='" + created_at + '\'' +
                ", updated_at='" + updated_at + '\'' +
                ", __op='" + __op + '\'' +
                ", __table='" + __table + '\'' +
                ", __lsn=" + __lsn +
                ", __source_ts_ms=" + __source_ts_ms +
                '}';
    }
}
