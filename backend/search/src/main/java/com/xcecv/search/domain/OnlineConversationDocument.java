package com.xcecv.search.domain;

import lombok.*;
import org.springframework.data.annotation.Id;
import org.springframework.data.elasticsearch.annotations.*;

import java.time.Instant;
import java.util.UUID;

@Document(indexName = "online_conversation")
@Setting(settingPath = "opensearch/standard-setting.json")
@Getter
@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class OnlineConversationDocument {
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

  @Field(type = FieldType.Date)
  private Instant time;
}
