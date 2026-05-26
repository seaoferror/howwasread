package com.xcecv.offlineconversation.domain;

import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.UuidGenerator;

import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

@Entity
@Builder
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Table(
    indexes = {
        @Index(name = "idx_h3_index", columnList = "h3_index")
    }
)
public class OfflineConversation {

  @Id
  @UuidGenerator(style = UuidGenerator.Style.VERSION_7)
  private UUID id;

  @Column(nullable = false)
  private String name;

  @Column(nullable = false)
  private double latitude;

  @Column(nullable = false)
  private double longitude;

  @Column(length = 15, nullable = false)
  private String h3Index;

  @ElementCollection(fetch = FetchType.LAZY)
  @CollectionTable(
      joinColumns = @JoinColumn(
          foreignKey = @ForeignKey(ConstraintMode.NO_CONSTRAINT)
      )
  )
  @Column(name = "participant_id", nullable = false)
  private Set<UUID> participants = new HashSet<>();
}
