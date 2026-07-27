package com.xcecv.offlineconversation.domain;

import jakarta.persistence.Id;
import lombok.*;
import org.springframework.data.elasticsearch.annotations.*;

import java.time.Instant;
import java.util.UUID;


@Document(indexName = "offline_conversation")
@Setting(settingPath = "opensearch/standard-setting.json")
@Getter
@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class OfflineConversationDocument {
  @Id
  private UUID id;

  @MultiField(
      mainField = @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer"),
      otherFields = {
          @InnerField(suffix = "keyword", type = FieldType.Keyword)
      }
  )
  private String novel;
  @MultiField(
      mainField = @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer"),
      otherFields = {
          @InnerField(suffix = "keyword", type = FieldType.Keyword)
      }
  )
  private String poem;
  @MultiField(
      mainField = @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer"),
      otherFields = {
          @InnerField(suffix = "keyword", type = FieldType.Keyword)
      }
  )
  private String shortStory;
  @MultiField(
      mainField = @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer"),
      otherFields = {
          @InnerField(suffix = "keyword", type = FieldType.Keyword)
      }
  )
  private String play;
  @MultiField(
      mainField = @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer"),
      otherFields = {
          @InnerField(suffix = "keyword", type = FieldType.Keyword)
      }
  )
  private String film;
  @MultiField(
      mainField = @Field(type = FieldType.Text, analyzer = "standard_lowercase_analyzer"),
      otherFields = {
          @InnerField(suffix = "keyword", type = FieldType.Keyword)
      }
  )
  private String writtenBy;

  @Field(type = FieldType.Date, format = {}, pattern = "uuuu-MM-dd'T'HH:mm:ss.SSSX")
  private Instant time;

  @Field(type = FieldType.Keyword)
  private String h3Res5;

  @Field(type = FieldType.Keyword)
  private String h3Res7;

  @Field(type = FieldType.Double, index = false)
  private Double latitude;

  @Field(type = FieldType.Double, index = false)
  private Double longitude;
}
