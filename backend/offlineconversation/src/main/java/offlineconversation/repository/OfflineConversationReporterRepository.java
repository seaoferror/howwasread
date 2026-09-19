package offlineconversation.repository;

import offlineconversation.domain.ConversationMemberCompositeKey;
import offlineconversation.domain.OfflineConversationReporter;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface OfflineConversationReporterRepository extends JpaRepository<OfflineConversationReporter, ConversationMemberCompositeKey> {
  long countByKeyConversationId(UUID conversationId);
}
