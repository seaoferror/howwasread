package com.howwasread;

import java.io.Serializable;

/**
 * Represents a notification event flowing through the Flink pipeline.
 * TODO: Add your actual fields and serialization logic.
 */
public class NotificationEvent implements Serializable {

    private static final long serialVersionUID = 1L;

    private String userId;
    private String message;
    private long timestamp;

    public NotificationEvent() {
    }

    public NotificationEvent(String userId, String message, long timestamp) {
        this.userId = userId;
        this.message = message;
        this.timestamp = timestamp;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public long getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(long timestamp) {
        this.timestamp = timestamp;
    }
}
