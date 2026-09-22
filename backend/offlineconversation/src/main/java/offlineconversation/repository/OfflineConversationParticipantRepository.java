package offlineconversation.repository;

import offlineconversation.domain.ConversationMemberCompositeKey;
import offlineconversation.domain.OfflineConversationParticipant;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;


@Repository
public interface OfflineConversationParticipantRepository extends JpaRepository<OfflineConversationParticipant, ConversationMemberCompositeKey> {
}
