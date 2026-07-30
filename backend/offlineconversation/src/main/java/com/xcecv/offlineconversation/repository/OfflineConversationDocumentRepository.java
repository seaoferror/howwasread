package com.xcecv.offlineconversation.repository;

import com.xcecv.offlineconversation.domain.OfflineConversationDocument;
import org.springframework.data.domain.Pageable;
import org.springframework.data.elasticsearch.annotations.Highlight;
import org.springframework.data.elasticsearch.annotations.HighlightField;
import org.springframework.data.elasticsearch.annotations.HighlightParameters;
import org.springframework.data.elasticsearch.annotations.Query;
import org.springframework.data.elasticsearch.core.SearchHit;
import org.springframework.data.elasticsearch.repository.ElasticsearchRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface OfflineConversationDocumentRepository extends ElasticsearchRepository<OfflineConversationDocument, UUID> {

  @Highlight(
      fields = {
          @HighlightField(name = "novel"),
          @HighlightField(name = "poem"),
          @HighlightField(name = "shortStory"),
          @HighlightField(name = "play"),
          @HighlightField(name = "film"),
          @HighlightField(name = "writtenBy")
      },
      parameters = @HighlightParameters(
          preTags = "<em>",
          postTags = "</em>"
      )
  )
  @Query("""
        {
          "bool": {
            "must": [
              {
                "multi_match": {
                  "query": "?0",
                  "fields": [
                    "novel", "novel.keyword^5",
                    "poem", "poem.keyword^5",
                    "shortStory", "shortStory.keyword^5",
                    "play", "play.keyword^5",
                    "film", "film.keyword^5",
                    "writtenBy", "writtenBy.keyword^5"
                  ],
                  "fuzziness": "AUTO"
                }
              }
            ],
            "filter": [
              {
                "terms": {
                  "h3Res7": ?1
                }
              },
              {
                "range": {
                  "time": {
                    "gt": "?2"
                  }
                }
              }
            ]
          }
        }
      """)
  List<SearchHit<OfflineConversationDocument>> findByInputAndH3Res7(
      String input,
      List<String> h3Indexes,
      String time,
      Pageable page
  );

  @Highlight(
      fields = {
          @HighlightField(name = "novel"),
          @HighlightField(name = "poem"),
          @HighlightField(name = "shortStory"),
          @HighlightField(name = "play"),
          @HighlightField(name = "film"),
          @HighlightField(name = "writtenBy")
      },
      parameters = @HighlightParameters(
          preTags = "<em>",
          postTags = "</em>"
      )
  )
  @Query("""
        {
          "bool": {
            "must": [
              {
                "multi_match": {
                  "query": "?0",
                  "fields": [
                    "novel", "novel.keyword^5",
                    "poem", "poem.keyword^5",
                    "shortStory", "shortStory.keyword^5",
                    "play", "play.keyword^5",
                    "film", "film.keyword^5",
                    "writtenBy", "writtenBy.keyword^5"
                  ],
                  "fuzziness": "AUTO"
                }
              }
            ],
            "filter": [
              {
                "terms": {
                  "h3Res5": ?1
                }
              },
              {
                "range": {
                  "time": {
                    "gt": "?2"
                  }
                }
              }
            ]
          }
        }
      """)
  List<SearchHit<OfflineConversationDocument>> findByInputAndH3Res5(
      String input,
      List<String> h3Indexes,
      String time,
      Pageable page
  );
}