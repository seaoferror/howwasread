package offlineconversation.domain;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.hibernate.annotations.CreationTimestamp;

import java.time.Instant;

@Entity
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class OfflineConversationReporter {
  @NotNull
  @EmbeddedId
  private ConversationMemberCompositeKey key;

  @CreationTimestamp
  @Column(updatable = false)
  private Instant createdAt;

  @NotNull
  @ManyToOne(fetch = FetchType.LAZY)
  @MapsId("conversationId")
  @JoinColumn(
      name = "conversation_id",
      foreignKey = @ForeignKey(ConstraintMode.NO_CONSTRAINT)
  )
  private OfflineConversation offlineConversation;
}
