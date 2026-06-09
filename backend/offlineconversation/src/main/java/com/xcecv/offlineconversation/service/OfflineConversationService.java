package com.xcecv.offlineconversation.service;

import com.xcecv.offlineconversation.domain.OfflineConversation;
import com.xcecv.offlineconversation.dto.CreateOfflineConversationRequest;
import com.xcecv.offlineconversation.dto.JoinOfflineConversationRequest;
import com.xcecv.offlineconversation.dto.OfflineConversationDetailResponse;
import com.xcecv.offlineconversation.dto.OfflineConversationMapResponse;
import com.xcecv.offlineconversation.projection.OfflineConversationDetailProjection;
import com.xcecv.offlineconversation.projection.OfflineConversationMapProjection;
import com.xcecv.offlineconversation.repository.OfflineConversationRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.time.Duration;
import java.util.*;

@Slf4j
@Service
@Transactional(readOnly = true)
@RequiredArgsConstructor
public class OfflineConversationService {

  private final OfflineConversationRepository offlineConversationRepository;

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
        .googleMapsLink(request.googleMapsLink())
        .location(request.location())
        .latitude(request.lat())
        .longitude(request.lng())
        .city(request.city())
        .h3Res5(request.h3Res5())
        .h3Res7(request.h3Res7())
        .moderatorIds(new HashSet<>(Set.of(memberId)))
        .participants(new HashSet<>(Set.of(memberId)))
        .build();

    return Map.of("id", offlineConversationRepository.save(convo).getId());
  }

  @Transactional
  public void join(JoinOfflineConversationRequest request, UUID memberId) {
    var convo = offlineConversationRepository.findById(request.conversationId())
        .orElseThrow(() -> new ResponseStatusException(
            HttpStatus.NOT_FOUND,
            "Conversation not found"
        ));

    convo.getParticipants().add(memberId);
  }

  public List<OfflineConversationMapResponse> mapRes7Convos(
      String h3Res7) {
    var convos = offlineConversationRepository.findByH3Res7(h3Res7);
    return buildOfflineConversationMapResponse(convos);
  }

  public List<OfflineConversationMapResponse> mapRes5Convos(
      String h3Res5) {
    var convos = offlineConversationRepository.findTop2ByH3Res5(h3Res5);
    return buildOfflineConversationMapResponse(convos);
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

  public OfflineConversationDetailResponse mapDetail(UUID conversationId, UUID memberId) {
    var convo = offlineConversationRepository.findById(conversationId, OfflineConversationDetailProjection.class)
        .orElseThrow(() -> new ResponseStatusException(
            HttpStatus.NOT_FOUND,
            "Conversation not found"
        ));
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
        .googleMapsLink(convo.getGoogleMapsLink())
        .location(convo.getLocation())
        .isModerator(convo.getModeratorIds().contains(memberId))
        .isParticipant(convo.getParticipants().contains(memberId))
        .numberOfParticipants(convo.getParticipants().size())
        .build();
  }
}
