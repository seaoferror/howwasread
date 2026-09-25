package offlineconversation.repository;

import offlineconversation.domain.Outbox;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface OutboxRepository extends JpaRepository<Outbox, UUID> {

  @Modifying
  @Query("DELETE FROM Outbox o WHERE o.id = :id AND o.conversationId = :conversationId")
  void deleteByIdAndConversationId(@Param("id") UUID id, @Param("conversationId") UUID conversationId);
}
