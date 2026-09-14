package exponentialbackoffretryjob.schema;

import lombok.extern.slf4j.Slf4j;
import org.apache.flink.api.common.serialization.SerializationSchema;
import exponentialbackoffretryjob.dto.OutgoingEvent;

import java.io.Serial;

@Slf4j
public class OutgoingEventSerializationSchema implements SerializationSchema<OutgoingEvent> {
  @Serial
  private static final long serialVersionUID = 1L;

  @Override
  public byte[] serialize(OutgoingEvent event) {
    if (event == null) {
      return null;
    }
    try {
      return event.getValue();
    } catch (Exception e) {
      log.error("fail to serialize: {}", e.getMessage());
      return null;
    }
  }
}
