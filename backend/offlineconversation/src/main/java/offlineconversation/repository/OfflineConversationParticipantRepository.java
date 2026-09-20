package offlineconversation.repository;

import offlineconversation.domain.ConversationMemberCompositeKey;
import offlineconversation.domain.OfflineConversationParticipant;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface OfflineConversationParticipantRepository extends JpaRepository<OfflineConversationParticipant, ConversationMemberCompositeKey> {
  @Modifying
  @Query("UPDATE OfflineConversationParticipant p SET p.deletedAt = CURRENT_TIMESTAMP WHERE p.key = :key")
  void softDeleteById(@Param("key") ConversationMemberCompositeKey key);
}
