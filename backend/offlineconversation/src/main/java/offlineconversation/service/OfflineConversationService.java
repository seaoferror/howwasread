package offlineconversation.service;

import offlineconversation.domain.ConversationMemberCompositeKey;
import offlineconversation.domain.OfflineConversation;
import offlineconversation.domain.OfflineConversationModerator;
import offlineconversation.domain.OfflineConversationParticipant;
import offlineconversation.domain.OfflineConversationReporter;
import offlineconversation.dto.*;
import offlineconversation.projection.OfflineConversationDetailProjection;
import offlineconversation.repository.OfflineConversationModeratorRepository;
import offlineconversation.repository.OfflineConversationParticipantRepository;
import offlineconversation.repository.OfflineConversationReporterRepository;
import offlineconversation.repository.OfflineConversationRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.util.*;

@Slf4j
@Service
@Transactional(readOnly = true)
@RequiredArgsConstructor
public class OfflineConversationService {

  private final OfflineConversationRepository offlineConversationRepository;
  private final OfflineConversationParticipantRepository offlineConversationParticipantRepository;
  private final OfflineConversationModeratorRepository offlineConversationModeratorRepository;
  private final OfflineConversationReporterRepository offlineConversationReporterRepository;

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
        .lengthMinutes(request.lengthMinutes())
        .mapsLink(request.mapsLink())
        .location(request.location())
        .latitude(request.lat())
        .longitude(request.lng())
        .city(request.city())
        .h3Res5(request.h3Res5())
        .h3Res7(request.h3Res7())
        .build();
    var conversationId = offlineConversationRepository.save(convo).getId();
    var key = ConversationMemberCompositeKey.builder()
        .conversationId(conversationId)
        .memberId(memberId)
        .build();
    offlineConversationParticipantRepository.save(new OfflineConversationParticipant(key, convo));
    offlineConversationModeratorRepository.save(new OfflineConversationModerator(key, convo));
    //TODO: CDC producing kafka message for chat group room create
    //TODO: CDC producing search
    return Map.of("id", conversationId);
  }

  @Transactional
  public void join(JoinOfflineConversationRequest request, UUID memberId) {
    var conversationProxy = offlineConversationRepository.getReferenceById(request.conversationId());
    var key = ConversationMemberCompositeKey.builder()
        .conversationId(request.conversationId())
        .memberId(memberId)
        .build();
    offlineConversationParticipantRepository.save(new OfflineConversationParticipant(key, conversationProxy));
    //TODO: CDC producing kafka message for chat group room participate
  }

  @Transactional
  public void quit(JoinOfflineConversationRequest request, UUID memberId) {
    var key = ConversationMemberCompositeKey.builder()
        .conversationId(request.conversationId())
        .memberId(memberId)
        .build();
    offlineConversationParticipantRepository.deleteById(key);
    //TODO: CDC producing kafka message for chat group room quit
  }

  public OfflineConversationDetailResponse detail(UUID conversationId, UUID memberId) {
    var convo = offlineConversationRepository.findById(conversationId, OfflineConversationDetailProjection.class)
        .orElseThrow(() -> new ResponseStatusException(
            HttpStatus.NOT_FOUND,
            "Conversation not found"
        ));
    var participantIds = offlineConversationParticipantRepository.findParticipantIdsByConversationId(conversationId);
    var moderatorIds = offlineConversationModeratorRepository.findMemberIdsByConversationId(conversationId);
    return OfflineConversationDetailResponse.builder()
        .novel(convo.getNovel())
        .poem(convo.getPoem())
        .shortStory(convo.getShortStory())
        .play(convo.getPlay())
        .film(convo.getFilm())
        .writtenBy(convo.getWrittenBy())
        .rule(convo.getRule())
        .time(convo.getTime())
        .lengthMinutes(convo.getLengthMinutes())
        .mapsLink(convo.getMapsLink())
        .location(convo.getLocation())
        .isModerator(moderatorIds.contains(memberId))
        .isParticipant(participantIds.contains(memberId))
        .numberOfParticipants(participantIds.size())
        .moderatorIds(moderatorIds)
        .build();
  }

  //TODO: refactor with jdbc
  @Transactional
  public void report(UUID conversationId, UUID memberId) {
    var key = ConversationMemberCompositeKey.builder()
        .conversationId(conversationId)
        .memberId(memberId)
        .build();
    if (offlineConversationReporterRepository.existsById(key)) {
      return;
    }
    long reporterCount = offlineConversationReporterRepository.countByKeyConversationId(conversationId);
    if (reporterCount > 5) {
      offlineConversationRepository.deleteById(conversationId);
      return;
    }
    var conversationProxy = offlineConversationRepository.getReferenceById(conversationId);
    offlineConversationReporterRepository.save(new OfflineConversationReporter(key, conversationProxy));
  }
}
