package exponentialbackoffretryjob.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;


@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OutgoingEvent implements Serializable {

  @Serial
  private static final long serialVersionUID = 1L;

  private String topic;
  private String originalTopic;
  private String type;
  private byte[] rawReasons;
  private byte[] value;
}
