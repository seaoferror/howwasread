package com.xcecv.search.service;

import com.xcecv.search.domain.OfflineConversationDocument;
import com.xcecv.search.domain.OnlineConversationDocument;
import com.xcecv.search.dto.OfflineConversationSearchResponse;
import com.xcecv.search.dto.OnlineConversationSearchResponse;
import com.xcecv.search.repository.OfflineConversationDocumentRepository;
import com.xcecv.search.repository.OnlineConversationDocumentRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
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

@Slf4j
@RequiredArgsConstructor
@Service
public class ConversationSearchService {

  private final OfflineConversationDocumentRepository offlineConversationDocumentRepository;
  private final OnlineConversationDocumentRepository onlineConversationDocumentRepository;
  private final ObjectMapper objectMapper;

  @KafkaListener(topics = "search", groupId = "search")
  public void consume(@Payload String payloadRaw, @Header("type") String type) {
    log.info("Successfully consumed message-type: {}", type);
    if (type.equals("offlineconversation")) {
      OfflineConversationDocument doc = objectMapper.readValue(payloadRaw, OfflineConversationDocument.class);
      offlineConversationDocumentRepository.save(doc);
      return;
    }
    if (type.equals("onlineconversation")) {
      OnlineConversationDocument doc = objectMapper.readValue(payloadRaw, OnlineConversationDocument.class);
      onlineConversationDocumentRepository.save(doc);
      return;
    }
    System.err.println("Unknown event type received: " + type);
  }

  public List<OfflineConversationSearchResponse> searchOfflines(String input, String resolution, List<String> h3Indexes, Instant time, int page) {
    if (resolution.equals("5")) {
      var searchHits = offlineConversationDocumentRepository
          .findByInputAndH3Res5(input, h3Indexes, time.toEpochMilli(), PageRequest.of(page - 1, 5));
      return buildOfflineConversationSearchResponses(searchHits);
    }
    if (resolution.equals("7")) {
      var searchHits = offlineConversationDocumentRepository
          .findByInputAndH3Res7(input, h3Indexes, time.toEpochMilli(), PageRequest.of(page - 1, 5));
      return buildOfflineConversationSearchResponses(searchHits);
    }
    return null;
  }

  public List<OnlineConversationSearchResponse> searchOnlines(String input, Instant time, int page) {
    var searchHits = onlineConversationDocumentRepository
        .findByInput(input, time.toEpochMilli(), PageRequest.of(page - 1, 5));
    return buildOnlineConversationSearchResponses(searchHits);
  }

  private List<OnlineConversationSearchResponse> buildOnlineConversationSearchResponses(List<SearchHit<OnlineConversationDocument>> searchHits) {
    List<OnlineConversationSearchResponse> response = new ArrayList<>();
    for (var searchHit : searchHits) {
      var conversation = searchHit.getContent();
      Map<String, List<String>> highlightFields = searchHit.getHighlightFields();
      response.add(OnlineConversationSearchResponse.builder()
          .id(conversation.getId())
          .novel(getHighlightOrOriginal(highlightFields, "novel", conversation.getNovel()))
          .play(getHighlightOrOriginal(highlightFields, "play", conversation.getPlay()))
          .poem(getHighlightOrOriginal(highlightFields, "poem", conversation.getPoem()))
          .shortStory(getHighlightOrOriginal(highlightFields, "shortStory", conversation.getShortStory()))
          .film(getHighlightOrOriginal(highlightFields, "film", conversation.getFilm()))
          .writtenBy(getHighlightOrOriginal(highlightFields, "writtenBy", conversation.getWrittenBy()))
          .time(conversation.getTime())
          .build());
    }
    return response;
  }

  private List<OfflineConversationSearchResponse> buildOfflineConversationSearchResponses(List<SearchHit<OfflineConversationDocument>> searchHits) {
    List<OfflineConversationSearchResponse> response = new ArrayList<>();
    for (var searchHit : searchHits) {
      var conversation = searchHit.getContent();
      Map<String, List<String>> highlightFields = searchHit.getHighlightFields();
      response.add(OfflineConversationSearchResponse.builder()
          .id(conversation.getId())
          .novel(getHighlightOrOriginal(highlightFields, "novel", conversation.getNovel()))
          .play(getHighlightOrOriginal(highlightFields, "play", conversation.getPlay()))
          .poem(getHighlightOrOriginal(highlightFields, "poem", conversation.getPoem()))
          .shortStory(getHighlightOrOriginal(highlightFields, "shortStory", conversation.getShortStory()))
          .film(getHighlightOrOriginal(highlightFields, "film", conversation.getFilm()))
          .writtenBy(getHighlightOrOriginal(highlightFields, "writtenBy", conversation.getWrittenBy()))
          .time(conversation.getTime())
          .lat(conversation.getLatitude())
          .lng(conversation.getLongitude())
          .build());
    }
    return response;
  }

  private String getHighlightOrOriginal(Map<String, List<String>> highlights, String fieldName, String originalValue) {
    return highlights.containsKey(fieldName) ? highlights.get(fieldName).getFirst() : originalValue;
  }
}
