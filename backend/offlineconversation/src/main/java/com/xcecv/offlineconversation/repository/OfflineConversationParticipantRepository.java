package com.xcecv.offlineconversation.repository;

import com.xcecv.offlineconversation.domain.OfflineConversationParticipant;
import com.xcecv.offlineconversation.domain.ParticipantCompositeKey;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;
import java.util.UUID;

public interface OfflineConversationParticipantRepository extends JpaRepository<OfflineConversationParticipant, ParticipantCompositeKey> {
  @Query("SELECT p.key.participantId FROM OfflineConversationParticipant p WHERE p.key.conversationId = :conversationId")
  List<UUID> findParticipantIdsByConversationId(@Param("conversationId") UUID conversationId);
}
