package com.xcecv.search.service;

import com.xcecv.search.domain.OfflineConversationDocument;
import com.xcecv.search.dto.OfflineConversationSearchResponse;
import com.xcecv.search.repository.OfflineConversationDocumentRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.elasticsearch.core.SearchHit;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.messaging.handler.annotation.Header;
import org.springframework.stereotype.Service;
import tools.jackson.databind.ObjectMapper;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@RequiredArgsConstructor
@Service
public class SearchService {

  private final OfflineConversationDocumentRepository offlineConversationDocumentRepository;
  private final ObjectMapper objectMapper;

  @KafkaListener(topics = "search", groupId = "search")
  public void consume(@Payload Object payloadRaw, @Header("type") String type) {
    if (type.equals("offlineconversation")) {
      OfflineConversationDocument doc = objectMapper.convertValue(payloadRaw, OfflineConversationDocument.class);
      offlineConversationDocumentRepository.save(doc);
      return;
    }
    System.err.println("Unknown event type received: " + type);
  }

  public List<OfflineConversationSearchResponse> searchH3Res5(String input, List<String> h3Indexes, int page) {
    var searchHits = offlineConversationDocumentRepository.findByInputAndH3Res5(input, h3Indexes, Instant.now().toEpochMilli(), PageRequest.of(page - 1, 5));
    return buildOfflineConversationSearchResponses(searchHits);
  }

  public List<OfflineConversationSearchResponse> searchH3Res7(String input, List<String> h3Indexes, int page) {
    var searchHits = offlineConversationDocumentRepository.findByInputAndH3Res7(input, h3Indexes, Instant.now().toEpochMilli(), PageRequest.of(page - 1, 5));
    return buildOfflineConversationSearchResponses(searchHits);
  }

  private List<OfflineConversationSearchResponse> buildOfflineConversationSearchResponses(List<SearchHit<OfflineConversationDocument>> searchHits) {
    List<OfflineConversationSearchResponse> response = new ArrayList<>();
    for (var searchHit : searchHits) {
      var conversation = searchHit.getContent();
      Map<String, List<String>> highlightFields = searchHit.getHighlightFields();
      response.add(OfflineConversationSearchResponse.builder().id(conversation.getId()).novel(getHighlightOrOriginal(highlightFields, "novel", conversation.getNovel())).play(getHighlightOrOriginal(highlightFields, "play", conversation.getPlay())).poem(getHighlightOrOriginal(highlightFields, "poem", conversation.getPoem())).shortStory(getHighlightOrOriginal(highlightFields, "shortStory", conversation.getShortStory())).film(getHighlightOrOriginal(highlightFields, "film", conversation.getFilm())).writtenBy(getHighlightOrOriginal(highlightFields, "writtenBy", conversation.getWrittenBy())).time(conversation.getTime()).lat(conversation.getLatitude()).lng(conversation.getLongitude()).build());
    }
    return response;
  }

  private String getHighlightOrOriginal(Map<String, List<String>> highlights, String fieldName, String originalValue) {
    return highlights.containsKey(fieldName) ? highlights.get(fieldName).getFirst() : originalValue;
  }
}
