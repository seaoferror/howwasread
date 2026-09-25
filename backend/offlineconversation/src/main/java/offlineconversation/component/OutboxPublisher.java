package offlineconversation.component;

import com.github.f4b6a3.uuid.UuidCreator;
import lombok.RequiredArgsConstructor;
import offlineconversation.domain.Outbox;
import offlineconversation.dto.ChatMessage;
import offlineconversation.repository.OutboxRepository;
import offlineconversation.util.UUIDUtil;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;
import tools.jackson.databind.ObjectMapper;

import java.util.List;
import java.util.UUID;

@Component
@RequiredArgsConstructor
@Transactional(propagation = Propagation.MANDATORY)
public class OutboxPublisher {

  private static final String CHAT_MESSAGE_TOPIC = "chat-message";

  private final OutboxRepository outboxRepository;
  private final ObjectMapper objectMapper;

  public void publishChatMessage(UUID conversationId, UUID memberId, String contentType, List<String> contents) {
    UUID id = UuidCreator.getTimeOrderedEpoch();
    String payload = objectMapper.writeValueAsString(ChatMessage.builder()
        .id(UUIDUtil.uuidToBytes(id))
        .fromId(UUIDUtil.uuidToBytes(memberId))
        .toIdType("group")
        .toId(UUIDUtil.uuidToBytes(conversationId))
        .contentType(contentType)
        .contents(contents)
        .build());
    publish(id, conversationId, CHAT_MESSAGE_TOPIC, payload);
  }

  private void publish(UUID id, UUID conversationId, String topic, String payload) {
    outboxRepository.saveAndFlush(Outbox.builder()
        .id(id)
        .conversationId(conversationId)
        .topic(topic)
        .payload(payload)
        .build());
    outboxRepository.deleteByIdAndConversationId(id, conversationId);
  }
}
