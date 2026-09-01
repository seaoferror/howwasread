package com.howwasread.mapper;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.howwasread.dto.IncomingNotificationEvent;
import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.api.common.typeinfo.TypeInformation;
import java.io.IOException;
import java.io.Serial;

public class NotificationEventDeserializationSchema implements DeserializationSchema<IncomingNotificationEvent> {
  @Serial
  private static final long serialVersionUID = 1L;

  private transient ObjectMapper objectMapper;

  @Override
  public void open(InitializationContext context) {
    objectMapper = new ObjectMapper();
  }

  @Override
  public IncomingNotificationEvent deserialize(byte[] message) throws IOException {
    if (message == null || message.length == 0) {
      return null;
    }
    return objectMapper.readValue(message, IncomingNotificationEvent.class);
  }

  @Override
  public boolean isEndOfStream(IncomingNotificationEvent nextElement) {
    return false;
  }

  @Override
  public TypeInformation<IncomingNotificationEvent> getProducedType() {
    return TypeInformation.of(IncomingNotificationEvent.class);
  }
}