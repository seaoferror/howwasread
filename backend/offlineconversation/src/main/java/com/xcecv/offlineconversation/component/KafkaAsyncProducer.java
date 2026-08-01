package com.xcecv.offlineconversation.component;

import com.xcecv.offlineconversation.dto.ChatMessage;
import com.xcecv.offlineconversation.dto.OfflineConversationDocument;
import com.xcecv.offlineconversation.util.UUIDUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.KafkaHeaders;
import org.springframework.messaging.Message;
import org.springframework.messaging.support.MessageBuilder;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Component;
import org.springframework.transaction.event.TransactionPhase;
import org.springframework.transaction.event.TransactionalEventListener;



@Slf4j
@Component
@RequiredArgsConstructor
public class KafkaAsyncProducer {

  private final KafkaTemplate<byte[], Object> kafkaTemplate;

  @Async
  @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
  public void handleParticipantUpdate(ChatMessage message) {
    kafkaTemplate.send("chat-message", message).whenComplete((result, ex) -> {
      if (ex != null) {
        log.error("fail to produce chat, ex, {}, id, {}, fromId, {}, toId, {}, contentType, {}, contents, {}",
            ex,
            UUIDUtil.bytesToHex(message.id()),
            UUIDUtil.bytesToHex(message.fromId()),
            UUIDUtil.bytesToHex(message.toId()),
            message.contentType(),
            message.contents());
      }
    });
  }

  @Async
  @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
  public void handleCreateConversation(OfflineConversationDocument document) {
    Message<OfflineConversationDocument> message = MessageBuilder
        .withPayload(document)
        .setHeader(KafkaHeaders.TOPIC, "search")
        .setHeader("type", "offlineconversation")
        .build();
    kafkaTemplate.send(message);
  }
}
