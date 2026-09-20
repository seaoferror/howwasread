package offlineconversation.repository;

import offlineconversation.domain.OfflineConversation;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface OfflineConversationRepository extends JpaRepository<OfflineConversation, UUID> {
  <T> Optional<T> findById(UUID id, Class<T> type);
}
