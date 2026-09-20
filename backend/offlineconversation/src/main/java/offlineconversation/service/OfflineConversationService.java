package offlineconversation.service;

import offlineconversation.domain.ConversationMemberCompositeKey;
import offlineconversation.domain.OfflineConversation;
import offlineconversation.domain.OfflineConversationModerator;
import offlineconversation.domain.OfflineConversationParticipant;
import offlineconversation.domain.OfflineConversationReporter;
import offlineconversation.dto.*;
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
      CreateOfflineConversationRequest req,
      UUID memberId
  ) {
    var convo = OfflineConversation.builder()
        .novel(req.novel())
        .poem(req.poem())
        .shortStory(req.shortStory())
        .play(req.play())
        .film(req.film())
        .writtenBy(req.writtenBy())
        .rule(req.rule())
        .time(req.time())
        .lengthMinutes(req.lengthMinutes())
        .mapsLink(req.mapsLink())
        .location(req.location())
        .latitude(req.lat())
        .longitude(req.lng())
        .city(req.city())
        .h3Res5(req.h3Res5())
        .h3Res7(req.h3Res7())
        .build();
    var conversationId = offlineConversationRepository.save(convo).getId();
    var key = ConversationMemberCompositeKey.builder()
        .conversationId(conversationId)
        .memberId(memberId)
        .build();
    offlineConversationParticipantRepository.save(
        new OfflineConversationParticipant(key, null, null, convo));
    offlineConversationModeratorRepository.save(
        new OfflineConversationModerator(key, null, null, convo));
    //TODO: CDC producing kafka message for chat group room create
    //TODO: CDC producing search
    return Map.of("id", conversationId);
  }

  @Transactional
  public void join(UUID conversationId, UUID memberId) {
    var conversationProxy = offlineConversationRepository.getReferenceById(conversationId);
    var key = ConversationMemberCompositeKey.builder()
        .conversationId(conversationId)
        .memberId(memberId)
        .build();
    offlineConversationParticipantRepository.save(
        new OfflineConversationParticipant(key, null, null, conversationProxy));
    //TODO: CDC producing kafka message for chat group room participate
  }

  @Transactional
  public void quit(UUID conversationId, UUID memberId) {
    var key = ConversationMemberCompositeKey.builder()
        .conversationId(conversationId)
        .memberId(memberId)
        .build();
    offlineConversationParticipantRepository.softDeleteById(key);
    //TODO: CDC producing kafka message for chat group room quit
  }

  public OfflineConversationDetailResponse detail(UUID conversationId, UUID memberId) {
    var convo = offlineConversationRepository.findDetail(conversationId, memberId)
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
        .lengthMinutes(convo.getLengthMinutes())
        .mapsLink(convo.getMapsLink())
        .location(convo.getLocation())
        .isModerator(convo.getIsModerator())
        .isParticipant(convo.getIsParticipant())
        .numberOfParticipants(convo.getNumberOfParticipants())
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
      offlineConversationParticipantRepository.softDeleteById(key);
      return;
    }
    var conversationProxy = offlineConversationRepository.getReferenceById(conversationId);
    offlineConversationReporterRepository.save(
        new OfflineConversationReporter(key, null, conversationProxy));
  }
}
