package offlineconversation.domain;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import lombok.*;

import java.io.Serializable;
import java.util.UUID;

@Builder
@Embeddable
@Getter
@Setter
@EqualsAndHashCode
@NoArgsConstructor
@AllArgsConstructor
public class ConversationMemberCompositeKey implements Serializable {
  @Column(nullable = false)
  private UUID conversationId;

  @Column(nullable = false)
  private UUID memberId;
}