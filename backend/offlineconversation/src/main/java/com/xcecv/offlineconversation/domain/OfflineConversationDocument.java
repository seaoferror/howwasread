package com.xcecv.offlineconversation.domain;

import jakarta.persistence.Id;
import lombok.*;
import org.springframework.data.elasticsearch.annotations.Document;
import org.springframework.data.elasticsearch.annotations.Field;
import org.springframework.data.elasticsearch.annotations.FieldType;
import org.springframework.data.elasticsearch.annotations.Setting;

import java.time.Instant;
import java.util.UUID;


@Document(indexName = "offline_conversation")
@Setting(settingPath = "standard-settings.json")
@Getter
@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class OfflineConversationDocument {
  @Id
  private UUID id;

  @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer")
  private String novel;

  @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer")
  private String poem;

  @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer")
  private String shortStory;

  @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer")
  private String play;

  @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer")
  private String film;

  @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer")
  private String writtenBy;

  @Field(type = FieldType.Date, format = {}, pattern = "uuuu-MM-dd'T'HH:mm:ss.SSSX")
  private Instant time;

  @Field(type = FieldType.Keyword)
  private String h3Res5;

  @Field(type = FieldType.Keyword)
  private String h3Res7;

  @Field(type = FieldType.Text, index = false)
  private String rule;

  @Field(type = FieldType.Keyword, index = false)
  private String mapsLink;

  @Field(type = FieldType.Text, index = false)
  private String location;

  @Field(type = FieldType.Double, index = false)
  private Double latitude;

  @Field(type = FieldType.Double, index = false)
  private Double longitude;
}
