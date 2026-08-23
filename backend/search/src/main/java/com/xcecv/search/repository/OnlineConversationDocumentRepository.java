package com.xcecv.search.repository;

import com.xcecv.search.domain.OfflineConversationDocument;
import com.xcecv.search.domain.OnlineConversationDocument;
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
public interface OnlineConversationDocumentRepository extends ElasticsearchRepository<OnlineConversationDocument, UUID> {

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
                  "type": "bool_prefix",
                  "fields": [
                    "novel", "novel._2gram", "novel._3gram",
                    "poem", "poem._2gram", "poem._3gram",
                    "shortStory", "shortStory._2gram", "shortStory._3gram",
                    "play", "play._2gram", "play._3gram",
                    "film", "film._2gram", "film._3gram",
                    "writtenBy", "writtenBy._2gram", "writtenBy._3gram"
                  ]
                }
              }
            ],
            "filter": [
              {
                "range": {
                  "time": {
                    "gt": ?1
                  }
                }
              }
            ]
          }
        }
      """)
  List<SearchHit<OnlineConversationDocument>> findByInput(
      String input,
      long time,
      Pageable page
  );
}
