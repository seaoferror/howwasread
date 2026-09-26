package offlineconversation.service;

import offlineconversation.component.OutboxPublisher;
import offlineconversation.domain.ConversationMemberCompositeKey;
import offlineconversation.domain.OfflineConversation;
import offlineconversation.domain.OfflineConversationModerator;
import offlineconversation.domain.OfflineConversationParticipant;
import offlineconversation.dto.CreateOfflineConversationRequest;
import offlineconversation.dto.OfflineConversationDetailResponse;
import offlineconversation.dto.UpdateOfflineConversationRequest;
import offlineconversation.projection.OfflineConversationDetailProjection;
import offlineconversation.repository.OfflineConversationModeratorRepository;
import offlineconversation.repository.OfflineConversationParticipantRepository;
import offlineconversation.repository.OfflineConversationRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.http.HttpStatus;
import org.springframework.web.server.ResponseStatusException;

import java.time.Instant;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatCode;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class OfflineConversationServiceTest {

  @Mock
  private OfflineConversationRepository offlineConversationRepository;
  @Mock
  private OfflineConversationParticipantRepository offlineConversationParticipantRepository;
  @Mock
  private OfflineConversationModeratorRepository offlineConversationModeratorRepository;
  @Mock
  private OutboxPublisher outboxPublisher;

  @InjectMocks
  private OfflineConversationService offlineConversationService;

  private final UUID conversationId = UUID.randomUUID();
  private final UUID memberId = UUID.randomUUID();
  private final Instant time = Instant.parse("2026-09-24T03:00:00Z");

  @Test
  void create_savesConversationMembersAndPublishesCreate() {
    when(offlineConversationRepository.save(any(OfflineConversation.class))).thenAnswer(invocation -> {
      OfflineConversation convo = invocation.getArgument(0);
      convo.setId(conversationId);
      return convo;
    });

    Map<String, UUID> response = offlineConversationService.create(createRequest(), memberId);

    assertThat(response).containsEntry("id", conversationId);
    ArgumentCaptor<OfflineConversation> convo = ArgumentCaptor.forClass(OfflineConversation.class);
    verify(offlineConversationRepository).save(convo.capture());
    assertThat(convo.getValue().getLocation()).isEqualTo("Seoul");
    assertThat(convo.getValue().getWrittenBy()).isEqualTo("shakespeare");
    assertThat(convo.getValue().getTime()).isEqualTo(time);

    ArgumentCaptor<OfflineConversationParticipant> participant = ArgumentCaptor.forClass(OfflineConversationParticipant.class);
    verify(offlineConversationParticipantRepository).save(participant.capture());
    assertThat(participant.getValue().getKey()).isEqualTo(key());

    ArgumentCaptor<OfflineConversationModerator> moderator = ArgumentCaptor.forClass(OfflineConversationModerator.class);
    verify(offlineConversationModeratorRepository).save(moderator.capture());
    assertThat(moderator.getValue().getKey()).isEqualTo(key());

    verify(outboxPublisher).publishChatMessage(conversationId, memberId, "create", List.of("Seoul"));
  }

  @Test
  void join_savesParticipantAndPublishesParticipate() {
    OfflineConversation proxy = OfflineConversation.builder().id(conversationId).build();
    when(offlineConversationRepository.getReferenceById(conversationId)).thenReturn(proxy);

    offlineConversationService.join(conversationId, memberId);

    ArgumentCaptor<OfflineConversationParticipant> participant = ArgumentCaptor.forClass(OfflineConversationParticipant.class);
    verify(offlineConversationParticipantRepository).save(participant.capture());
    assertThat(participant.getValue().getKey()).isEqualTo(key());
    assertThat(participant.getValue().getOfflineConversation()).isSameAs(proxy);
    verify(outboxPublisher).publishChatMessage(conversationId, memberId, "participate", List.of());
  }

  @Test
  void quit_deletesParticipantAndPublishesQuit() {
    offlineConversationService.quit(conversationId, memberId);

    verify(offlineConversationParticipantRepository).deleteById(key());
    verify(outboxPublisher).publishChatMessage(conversationId, memberId, "quit", List.of());
  }

  @Test
  void delete_succeedsWhenRowDeleted() {
    when(offlineConversationRepository.deleteIfModerator(conversationId, memberId)).thenReturn(1);

    assertThatCode(() -> offlineConversationService.delete(conversationId, memberId)).doesNotThrowAnyException();
  }

  @Test
  void delete_throwsBadRequestWhenNotModerator() {
    when(offlineConversationRepository.deleteIfModerator(conversationId, memberId)).thenReturn(0);

    assertThatThrownBy(() -> offlineConversationService.delete(conversationId, memberId))
        .isInstanceOfSatisfying(ResponseStatusException.class, e -> {
          assertThat(e.getStatusCode()).isEqualTo(HttpStatus.BAD_REQUEST);
          assertThat(e.getReason()).isEqualTo("can't delete conversation");
        });
  }

  @Test
  void update_succeedsWhenRowUpdated() {
    UpdateOfflineConversationRequest req = updateRequest();
    when(offlineConversationRepository.updateIfModerator(req, memberId)).thenReturn(1);

    assertThatCode(() -> offlineConversationService.update(req, memberId)).doesNotThrowAnyException();
  }

  @Test
  void update_throwsBadRequestWhenNotModerator() {
    UpdateOfflineConversationRequest req = updateRequest();
    when(offlineConversationRepository.updateIfModerator(req, memberId)).thenReturn(0);

    assertThatThrownBy(() -> offlineConversationService.update(req, memberId))
        .isInstanceOfSatisfying(ResponseStatusException.class, e -> {
          assertThat(e.getStatusCode()).isEqualTo(HttpStatus.BAD_REQUEST);
          assertThat(e.getReason()).isEqualTo("can't update conversation");
        });
  }

  @Test
  void detail_mapsProjection() {
    OfflineConversationDetailProjection projection = mock(OfflineConversationDetailProjection.class);
    when(projection.getNovel()).thenReturn("Hamlet");
    when(projection.getWrittenBy()).thenReturn("shakespeare");
    when(projection.getTime()).thenReturn(time);
    when(projection.getLengthMinutes()).thenReturn(60);
    when(projection.getMapsLink()).thenReturn("https://maps");
    when(projection.getLocation()).thenReturn("Seoul");
    when(projection.getIsModerator()).thenReturn(1L);
    when(projection.getIsParticipant()).thenReturn(0L);
    when(projection.getNumberOfParticipants()).thenReturn(3);
    when(offlineConversationRepository.findDetail(conversationId, memberId)).thenReturn(Optional.of(projection));

    OfflineConversationDetailResponse response = offlineConversationService.detail(conversationId, memberId);

    assertThat(response.novel()).isEqualTo("Hamlet");
    assertThat(response.writtenBy()).isEqualTo("shakespeare");
    assertThat(response.time()).isEqualTo(time);
    assertThat(response.lengthMinutes()).isEqualTo(60);
    assertThat(response.mapsLink()).isEqualTo("https://maps");
    assertThat(response.location()).isEqualTo("Seoul");
    assertThat(response.isModerator()).isTrue();
    assertThat(response.isParticipant()).isFalse();
    assertThat(response.numberOfParticipants()).isEqualTo(3);
  }

  @Test
  void detail_throwsNotFoundWhenMissing() {
    when(offlineConversationRepository.findDetail(conversationId, memberId)).thenReturn(Optional.empty());

    assertThatThrownBy(() -> offlineConversationService.detail(conversationId, memberId))
        .isInstanceOfSatisfying(ResponseStatusException.class,
            e -> assertThat(e.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND));
    verifyNoInteractions(outboxPublisher);
  }

  private ConversationMemberCompositeKey key() {
    return ConversationMemberCompositeKey.builder()
        .conversationId(conversationId)
        .memberId(memberId)
        .build();
  }

  private CreateOfflineConversationRequest createRequest() {
    return new CreateOfflineConversationRequest("Hamlet", null, null, null, null, "shakespeare", null,
        time, 60, "https://maps", "Seoul", "Seoul", 37.5, 127.0, "85283473fffffff", "87283472bffffff");
  }

  private UpdateOfflineConversationRequest updateRequest() {
    return new UpdateOfflineConversationRequest(conversationId, "Hamlet", null, null, null, null, "shakespeare",
        null, time, 60, "https://maps", "Seoul", "Seoul", 37.5, 127.0, "85283473fffffff", "87283472bffffff");
  }
}
