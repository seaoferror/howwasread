package com.xcecv.offlineconversation.service;

import com.xcecv.offlineconversation.domain.OfflineConversation;
import com.xcecv.offlineconversation.dto.CreateOfflineConversationRequest;
import com.xcecv.offlineconversation.dto.JoinOfflineConversationRequest;
import com.xcecv.offlineconversation.dto.ShowOfflineConversationResponse;
import com.xcecv.offlineconversation.repository.OfflineConversationRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.util.*;

@Service
@Transactional(readOnly = true)
@RequiredArgsConstructor
public class OfflineConversationService {

  private final OfflineConversationRepository repository;

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
        .by(request.by())
        .rule(request.rule())
        .capacity(request.capacity())
        .when(request.when())
        .where(request.where())
        .latitude(request.lat())
        .longitude(request.lng())
        .city(request.city())
        .h3Res7(request.h3Index())
        .moderatorIds(new HashSet<>(Set.of(memberId)))
        .participants(new HashSet<>(Set.of(memberId)))
        .build();

    return Map.of("id", repository.save(convo).getId());
  }

  @Transactional
  public void join(JoinOfflineConversationRequest request, UUID memberId) {
    var convo = repository.findById(request.conversationId())
        .orElseThrow(() -> new ResponseStatusException(
            HttpStatus.NOT_FOUND,
            "Conversation not found"
        ));

    convo.getParticipants().add(memberId);
  }

  public List<ShowOfflineConversationResponse> showConvos(
      String h3Index,
      UUID memberId) {
    var convos = repository.findByH3Res7(h3Index);
    List<ShowOfflineConversationResponse> response = new ArrayList<>();
    for (var convo : convos) {
      response.add(ShowOfflineConversationResponse.builder()
          .id(convo.getId())
          .novel(convo.getNovel())
          .poem(convo.getPoem())
          .shortStory(convo.getShortStory())
          .play(convo.getPlay())
          .film(convo.getFilm())
          .by(convo.getBy())
          .rule(convo.getRule())
          .capacity(convo.getCapacity())
          .when(convo.getWhen())
          .where(convo.getWhere())
          .lat(convo.getLatitude())
          .lng(convo.getLongitude())
          .isModerator(convo.getModeratorIds().contains(memberId))
          .isParticipant(convo.getParticipants().contains(memberId))
          .build());
    }
    return response;
  }
}
