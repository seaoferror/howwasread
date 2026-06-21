package com.xcecv.offlineconversation.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.f4b6a3.uuid.UuidCreator;
import com.xcecv.offlineconversation.domain.OfflineConversation;
import com.xcecv.offlineconversation.domain.OfflineConversationParticipant;
import com.xcecv.offlineconversation.domain.ParticipantCompositeKey;
import com.xcecv.offlineconversation.dto.*;
import com.xcecv.offlineconversation.projection.OfflineConversationDetailProjection;
import com.xcecv.offlineconversation.projection.OfflineConversationMapProjection;
import com.xcecv.offlineconversation.projection.OfflineConversationPinProjection;
import com.xcecv.offlineconversation.repository.OfflineConversationParticipantRepository;
import com.xcecv.offlineconversation.repository.OfflineConversationRepository;
import com.xcecv.offlineconversation.util.UUIDUtil;
import glide.api.GlideClusterClient;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;

@Slf4j
@Service
@Transactional(readOnly = true)
@RequiredArgsConstructor
public class OfflineConversationService {

  private final OfflineConversationRepository offlineConversationRepository;
  private final OfflineConversationParticipantRepository offlineConversationParticipantRepository;
  private final GlideClusterClient glideClusterClient;
  private final ObjectMapper objectMapper;
  private final KafkaTemplate<byte[], Object> kafkaTemplate;

  private static final String EMPTY_CACHE_DUMMY_KEY = "_";
  private static final long CACHE_TTL_SECONDS = 3600;

  @Transactional
  public Map<String, UUID> create(
      CreateOfflineConversationRequest request,
      UUID memberId
  ) {
    var convo = OfflineConversation.builder()
        .novel(request.novel())
        .poem(request.poem())
        .shortStory(request.shortStory())
        .play(request.play())
        .film(request.film())
        .writtenBy(request.writtenBy())
        .rule(request.rule())
        .time(request.time())
        .length(Duration.ofMinutes(request.length()))
        .mapsLink(request.mapsLink())
        .location(request.location())
        .latitude(request.lat())
        .longitude(request.lng())
        .city(request.city())
        .h3Res5(request.h3Res5())
        .h3Res7(request.h3Res7())
        .moderatorIds(new HashSet<>(Set.of(memberId)))
        .build();
    var conversationId = offlineConversationRepository.save(convo).getId();
    var key = ParticipantCompositeKey.builder()
        .conversationId(conversationId)
        .participantId(memberId)
        .build();
    offlineConversationParticipantRepository.save(new OfflineConversationParticipant(key, convo));
    kafkaTemplate.send("chat-message",
        ChatMessage.builder()
            .id(UUIDUtil.uuidToBytes(UuidCreator.getTimeOrderedEpoch()))
            .fromId(UUIDUtil.uuidToBytes(memberId))
            .toIdType("group")
            .toId(UUIDUtil.uuidToBytes(conversationId))
            .contentType("create")
            .contents(new ArrayList<>(List.of(request.location())))
            .build());
    try {
      Map<String, String> hashEntry = Map.of(conversationId.toString(),
          objectMapper.writeValueAsString(
              OfflineConversationPinProjection.builder()
                  .writtenBy(request.writtenBy())
                  .latitude(request.lat())
                  .longitude(request.lng())
                  .build()));
      var checkRes5 = glideClusterClient.exists(new String[]{request.h3Res5()});
      var checkRes7 = glideClusterClient.exists(new String[]{request.h3Res7()});
      List<CompletableFuture<?>> tasks = new ArrayList<>();
      if (checkRes5.join() > 0) {
        tasks.add(glideClusterClient.hset(request.h3Res5(), hashEntry));
        tasks.add(glideClusterClient.expire(request.h3Res5(), CACHE_TTL_SECONDS));
      }
      if (checkRes7.join() > 0) {
        tasks.add(glideClusterClient.hset(request.h3Res7(), hashEntry));
        tasks.add(glideClusterClient.expire(request.h3Res7(), CACHE_TTL_SECONDS));
      }
      if (!tasks.isEmpty()) {
        CompletableFuture.allOf(tasks.toArray(new CompletableFuture[0])).join();
      }
    } catch (Exception e) {
      log.error("Failed to push conversation {} to Valkey H3 cache", conversationId, e);
    }
    return Map.of("id", conversationId);
  }

  @Transactional
  public void join(JoinOfflineConversationRequest request, UUID memberId) {
    var conversationProxy = offlineConversationRepository.getReferenceById(request.conversationId());
    var key = ParticipantCompositeKey.builder()
        .conversationId(request.conversationId())
        .participantId(memberId)
        .build();
    offlineConversationParticipantRepository.save(new OfflineConversationParticipant(key, conversationProxy));
    kafkaTemplate.send("chat-message",
        ChatMessage.builder()
            .id(UUIDUtil.uuidToBytes(UuidCreator.getTimeOrderedEpoch()))
            .fromId(UUIDUtil.uuidToBytes(memberId))
            .toIdType("group")
            .toId(UUIDUtil.uuidToBytes(request.conversationId()))
            .contentType("participate")
            .contents(null)
            .build());
  }

  @Transactional
  public void quit(JoinOfflineConversationRequest request, UUID memberId) {
    var key = ParticipantCompositeKey.builder()
        .conversationId(request.conversationId())
        .participantId(memberId)
        .build();
    offlineConversationParticipantRepository.deleteById(key);
    kafkaTemplate.send("chat-message",
        ChatMessage.builder()
            .id(UUIDUtil.uuidToBytes(UuidCreator.getTimeOrderedEpoch()))
            .fromId(UUIDUtil.uuidToBytes(memberId))
            .toIdType("group")
            .toId(UUIDUtil.uuidToBytes(request.conversationId()))
            .contentType("quit")
            .contents(null)
            .build());
  }

  public List<OfflineConversationMapResponse> mapRes7Convos(
      String h3Res7) {
    try {
      Map<String, String> cache = glideClusterClient.hgetall(h3Res7).join();
      if (cache != null && !cache.isEmpty()) {
        return unmarshalH3Cache(cache);
      }
      var convos = offlineConversationRepository.findByH3Res7(h3Res7);
      var response = buildOfflineConversationMapResponse(convos);
      setH3Cache(h3Res7, response);
      return response;
    } catch (Exception e) {
      log.error("Failed to read from cache for h3Index {}. Falling back to DB.", h3Res7, e);
      var convos = offlineConversationRepository.findByH3Res7(h3Res7);
      return buildOfflineConversationMapResponse(convos);
    }
  }

  public List<OfflineConversationMapResponse> mapRes5Convos(
      String h3Res5) {
    try {
      Map<String, String> cacheRaw = glideClusterClient.hgetall(h3Res5).join();
      if (cacheRaw != null && !cacheRaw.isEmpty()) {
        return unmarshalH3Cache(cacheRaw);
      }
      var convos = offlineConversationRepository.findTop2ByH3Res5(h3Res5);
      var response = buildOfflineConversationMapResponse(convos);
      setH3Cache(h3Res5, response);
      return response;
    } catch (Exception e) {
      log.error("Failed to read from cache for h3Index {}. Falling back to DB.", h3Res5, e);
      var convos = offlineConversationRepository.findTop2ByH3Res5(h3Res5);
      return buildOfflineConversationMapResponse(convos);
    }
  }


  private List<OfflineConversationMapResponse> unmarshalH3Cache(Map<String, String> cacheRaw) throws JsonProcessingException {
    log.info("cache hit");
    List<OfflineConversationMapResponse> response = new ArrayList<>();
    for (Map.Entry<String, String> entry : cacheRaw.entrySet()) {
      if (EMPTY_CACHE_DUMMY_KEY.equals(entry.getKey())) {
        continue;
      }
      UUID conversationId = UUID.fromString(entry.getKey());
      OfflineConversationPinProjection pinValue = objectMapper.readValue(
          entry.getValue(),
          OfflineConversationPinProjection.class
      );
      response.add(OfflineConversationMapResponse.builder()
          .id(conversationId)
          .writtenBy(pinValue.writtenBy())
          .lat(pinValue.latitude())
          .lng(pinValue.longitude())
          .build());
    }
    return response;
  }

  private void setH3Cache(String h3Index, List<OfflineConversationMapResponse> response) throws JsonProcessingException {
    Map<String, String> values = new HashMap<>();
    if (response.isEmpty()) {
      values.put(EMPTY_CACHE_DUMMY_KEY, "1");
    }
    for (var item : response) {
      values.put(item.id().toString(),
          objectMapper.writeValueAsString(
              OfflineConversationPinProjection.builder()
                  .writtenBy(item.writtenBy())
                  .latitude(item.lat())
                  .longitude(item.lng())
                  .build()
          ));
    }
    glideClusterClient.hset(h3Index, values).join();
    glideClusterClient.expire(h3Index, CACHE_TTL_SECONDS).join();
  }


  private List<OfflineConversationMapResponse> buildOfflineConversationMapResponse(List<OfflineConversationMapProjection> convos) {
    List<OfflineConversationMapResponse> response = new ArrayList<>();
    for (var convo : convos) {
      response.add(OfflineConversationMapResponse.builder()
          .id(convo.getId())
          .writtenBy(convo.getWrittenBy())
          .lat(convo.getLatitude())
          .lng(convo.getLongitude())
          .build());
    }
    return response;
  }

  public OfflineConversationDetailResponse detail(UUID conversationId, UUID memberId) {
    var convo = offlineConversationRepository.findById(conversationId, OfflineConversationDetailProjection.class)
        .orElseThrow(() -> new ResponseStatusException(
            HttpStatus.NOT_FOUND,
            "Conversation not found"
        ));
    var participantIds = offlineConversationParticipantRepository.findParticipantIdsByConversationId(conversationId);
    return OfflineConversationDetailResponse.builder()
        .novel(convo.getNovel())
        .poem(convo.getPoem())
        .shortStory(convo.getShortStory())
        .play(convo.getPlay())
        .film(convo.getFilm())
        .writtenBy(convo.getWrittenBy())
        .rule(convo.getRule())
        .time(convo.getTime())
        .length((int) convo.getLength().toMinutes())
        .mapsLink(convo.getMapsLink())
        .location(convo.getLocation())
        .isModerator(convo.getModeratorIds().contains(memberId))
        .isParticipant(participantIds.contains(memberId))
        .numberOfParticipants(participantIds.size())
        .build();
  }
}
