package offlineconversation.repository;

import offlineconversation.domain.OfflineConversationParticipant;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface OfflineConversationParticipantRepository extends JpaRepository<OfflineConversationParticipant, UUID> {
  @Modifying
  @Query("UPDATE OfflineConversationParticipant p " +
      "SET p.deletedAt = CURRENT_TIMESTAMP " +
      "WHERE p.offlineConversation.id = :conversationId " +
      "AND p.memberId = :memberId " +
      "AND p.deletedAt IS NULL")
  void softDeleteByConversationIdAndMemberId(@Param("conversationId") UUID conversationId, @Param("memberId") UUID memberId);

  @Modifying
  @Query(value = """
      INSERT INTO offline_conversation_participant
      (id, conversation_id, member_id)
      SELECT :id, :conversationId, :memberId WHERE NOT EXISTS
      (SELECT 1 FROM offline_conversation_participant
      WHERE conversation_id = :conversationId
      AND member_id = :memberId
      AND deleted_at IS NULL)
      """, nativeQuery = true)
  void insertByConversationIdAndMemberId(@Param("id") UUID id, @Param("conversationId") UUID conversationId, @Param("memberId") UUID memberId);
  //TODO: figure out how native query prevent sql injection
}
