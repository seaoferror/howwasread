package offlineconversation.repository;

import offlineconversation.domain.ConversationMemberCompositeKey;
import offlineconversation.domain.OfflineConversationModerator;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface OfflineConversationModeratorRepository extends JpaRepository<OfflineConversationModerator, ConversationMemberCompositeKey> {
  @Query("SELECT m.key.memberId FROM OfflineConversationModerator m WHERE m.key.conversationId = :conversationId")
  List<UUID> findMemberIdsByConversationId(@Param("conversationId") UUID conversationId);
}
