package offlineconversation.repository;

import offlineconversation.domain.ConversationMemberCompositeKey;
import offlineconversation.domain.OfflineConversationParticipant;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface OfflineConversationParticipantRepository extends JpaRepository<OfflineConversationParticipant, ConversationMemberCompositeKey> {
  @Query("SELECT p.key.memberId FROM OfflineConversationParticipant p WHERE p.key.conversationId = :conversationId")
  List<UUID> findParticipantIdsByConversationId(@Param("conversationId") UUID conversationId);
}
