package offlineconversation.domain;

import jakarta.persistence.*;
import lombok.*;
import org.springframework.data.domain.Persistable;

import java.util.UUID;

/**
 * Transactional outbox row, inserted and deleted in the same transaction.
 * The CDC job forwards the insert to the Kafka topic in the topic column.
 */
@Entity
@Builder
@Getter
@NoArgsConstructor
@AllArgsConstructor
public class Outbox implements Persistable<UUID> {
  @Id
  private UUID id;

  @Column(name = "conversation_id")
  private UUID conversationId;

  private String topic;

  @Column(columnDefinition = "TEXT")
  private String payload;

  // always a new row, so save() persists instead of merging with a select by id (not the shard key)
  @Override
  public boolean isNew() {
    return true;
  }
}
