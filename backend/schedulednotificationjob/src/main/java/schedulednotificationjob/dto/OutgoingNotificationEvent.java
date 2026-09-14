package schedulednotificationjob.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.Map;


@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OutgoingNotificationEvent implements Serializable {

  @Serial
  private static final long serialVersionUID = 1L;

  private String partitionId;
  private String partitionType;
  private long scheduledTime;
  private Map<Integer, String> sharedContents;
  private Map<String, Map<Integer, String>> notifications;
}
