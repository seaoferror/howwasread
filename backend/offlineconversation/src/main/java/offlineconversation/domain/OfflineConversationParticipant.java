package offlineconversation.domain;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import lombok.*;
import org.hibernate.annotations.UuidGenerator;

import java.time.Instant;
import java.util.UUID;

@Entity
@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OfflineConversationParticipant {
  @Id
  @UuidGenerator(style = UuidGenerator.Style.VERSION_7)
  private UUID id;

  @Column(nullable = false)
  private UUID memberId;

  @Column
  private Instant deletedAt;

  @NotNull
  @ManyToOne(fetch = FetchType.LAZY)
  @JoinColumn(
      name = "conversation_id",
      foreignKey = @ForeignKey(ConstraintMode.NO_CONSTRAINT)
  )
  private OfflineConversation offlineConversation;
}