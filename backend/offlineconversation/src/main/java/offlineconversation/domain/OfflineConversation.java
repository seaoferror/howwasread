package offlineconversation.domain;

import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.UuidGenerator;

import java.time.Instant;
import java.util.UUID;

@Entity
@Builder
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Table(
    indexes = {
        @Index(name = "idx_h3_res5_time", columnList = "h3_res5, time"),
        @Index(name = "idx_h3_res7_time", columnList = "h3_res7, time")
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
  private int lengthMinutes;

  @Column(columnDefinition = "TEXT", nullable = false)
  private String mapsLink;

  @Column(columnDefinition = "TEXT")
  private String location;

  @Column(nullable = false)
  private double latitude;

  @Column(nullable = false)
  private double longitude;

  @Column(columnDefinition = "TEXT")
  private String city;

  @Column(length = 15, nullable = false)
  private String h3Res5;

  @Column(length = 15, nullable = false)
  private String h3Res7;

  @Column
  private Instant deletedAt;
}
