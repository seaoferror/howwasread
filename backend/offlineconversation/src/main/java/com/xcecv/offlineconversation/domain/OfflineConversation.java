package com.xcecv.offlineconversation.domain;

import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.annotations.UuidGenerator;
import org.hibernate.type.SqlTypes;

import java.time.Duration;
import java.time.Instant;
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
        @Index(name = "h3_res5", columnList = "h3_res5"),
        @Index(name = "h3_res7", columnList = "h3_res7")
    }
)
public class OfflineConversation {
  @Id
  @UuidGenerator(style = UuidGenerator.Style.VERSION_7)
  private UUID id;

  @Column(columnDefinition = "TEXT")
  private String novel;

  @Column(columnDefinition = "TEXT")
  private String poem;

  @Column(columnDefinition = "TEXT")
  private String shortStory;

  @Column(columnDefinition = "TEXT")
  private String play;

  @Column(columnDefinition = "TEXT")
  private String film;

  @Column(columnDefinition = "TEXT", nullable = false)
  private String writtenBy;

  @Column(columnDefinition = "TEXT")
  private String rule;

  @Column(nullable = false)
  private Instant time;

  @Column(nullable = false)
  private Duration length;

  @Column(nullable = false)
  private String googleMapsLink;

  private String location;

  @Column(nullable = false)
  private double latitude;

  @Column(nullable = false)
  private double longitude;

  private String city;

  @Column(length = 15, nullable = false)
  private String h3Res5;

  @Column(length = 15, nullable = false)
  private String h3Res7;

  @JdbcTypeCode(SqlTypes.JSON)
  @Column(columnDefinition = "JSON")
  private Set<UUID> moderatorIds = new HashSet<>();

  @ElementCollection(fetch = FetchType.LAZY)
  @CollectionTable(
      name = "offline_conversation_participants",
      joinColumns = @JoinColumn(
          foreignKey = @ForeignKey(ConstraintMode.NO_CONSTRAINT)
      )
  )
  @Column(name = "participant_id", nullable = false)
  private Set<UUID> participants = new HashSet<>();
}
