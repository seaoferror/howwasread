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
                "bool": {
                  "should": [
                    {
                      "query_string": {
                        "query": "*?0*",
                        "fields": [
                          "novel",
                          "poem",
                          "shortStory",
                          "play",
                          "film",
                          "writtenBy"
                        ]
                      }
                    },
                    {
                      "multi_match": {
                        "query": "?0",
                        "fields": [
                          "novel^3",
                          "poem^3",
                          "shortStory^3",
                          "play^3",
                          "film^3",
                          "writtenBy^3"
                        ]
                      }
                    },
                    {
                      "multi_match": {
                        "query": "?0",
                        "fields": [
                          "novel.keyword^5",
                          "poem.keyword^5",
                          "shortStory.keyword^5",
                          "play.keyword^5",
                          "film.keyword^5",
                          "writtenBy.keyword^5"
                        ]
                      }
                    }
                  ],
                  "minimum_should_match": 1
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
