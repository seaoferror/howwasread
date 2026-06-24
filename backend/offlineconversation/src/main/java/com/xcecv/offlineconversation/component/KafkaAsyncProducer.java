package com.xcecv.offlineconversation.component;

import com.xcecv.offlineconversation.dto.ChatMessage;
import lombok.RequiredArgsConstructor;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Component;
import org.springframework.transaction.event.TransactionPhase;
import org.springframework.transaction.event.TransactionalEventListener;


@Component
@RequiredArgsConstructor
public class KafkaAsyncProducer {

  private final KafkaTemplate<byte[], Object> kafkaTemplate;

  @Async
  @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
  public void handleParticipantUpdate(ChatMessage message) {
    kafkaTemplate.send("chat-message", message);
  }
}
