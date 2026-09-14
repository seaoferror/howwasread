package exponentialbackoffretryjob.schema;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.api.common.typeinfo.TypeInformation;
import exponentialbackoffretryjob.dto.IncomingEvent;

import java.io.IOException;
import java.io.Serial;

public class IncomingEventDeserializationSchema implements DeserializationSchema<IncomingEvent> {
  @Serial
  private static final long serialVersionUID = 1L;

  private transient ObjectMapper objectMapper;

  @Override
  public void open(InitializationContext context) {
    objectMapper = new ObjectMapper();
  }

  @Override
  public IncomingEvent deserialize(byte[] message) throws IOException {
    if (message == null || message.length == 0) {
      return null;
    }
    return objectMapper.readValue(message, IncomingEvent.class);
  }

  @Override
  public boolean isEndOfStream(IncomingEvent nextElement) {
    return false;
  }

  @Override
  public TypeInformation<IncomingEvent> getProducedType() {
    return TypeInformation.of(IncomingEvent.class);
  }
}
