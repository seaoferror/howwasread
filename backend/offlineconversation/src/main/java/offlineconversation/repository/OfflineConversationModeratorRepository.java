package offlineconversation.repository;

import offlineconversation.domain.OfflineConversationModerator;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface OfflineConversationModeratorRepository extends JpaRepository<OfflineConversationModerator, UUID> {
}
