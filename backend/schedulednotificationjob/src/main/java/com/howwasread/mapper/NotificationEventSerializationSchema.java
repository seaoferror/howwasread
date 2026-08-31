package com.howwasread.mapper;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.howwasread.dto.OutgoingNotificationEvent;
import lombok.extern.slf4j.Slf4j;
import org.apache.flink.api.common.serialization.SerializationSchema;

import java.io.Serial;

@Slf4j
public class NotificationEventSerializationSchema implements SerializationSchema<OutgoingNotificationEvent> {
  @Serial
  private static final long serialVersionUID = 1L;
  private transient ObjectMapper objectMapper;

  @Override
  public void open(InitializationContext context) {
    objectMapper = new ObjectMapper();
  }

  @Override
  public byte[] serialize(OutgoingNotificationEvent element) {
    if (element == null) {
      return null;
    }
    try {
      return objectMapper.writeValueAsBytes(element);
    } catch (Exception e) {
      log.error("fail to serialize: {}", e.getMessage());
      return null;
    }
  }
}