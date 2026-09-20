package offlineconversation.repository;

import offlineconversation.domain.OfflineConversation;
import offlineconversation.projection.OfflineConversationDetailProjection;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface OfflineConversationRepository extends JpaRepository<OfflineConversation, UUID> {

  @Query(value = """
      SELECT 
        c.novel as novel, c.poem as poem, c.short_story as shortStory, 
        c.play as play, c.film as film, c.written_by as writtenBy, c.rule as rule, 
        c.time as time, c.length_minutes as lengthMinutes, c.maps_link as mapsLink, 
        c.location as location,
        EXISTS(SELECT 1 FROM offline_conversation_moderator m WHERE m.conversation_id = c.id AND m.member_id = :memberId AND m.deleted_at IS NULL) as isModerator,
        EXISTS(SELECT 1 FROM offline_conversation_participant p WHERE p.conversation_id = c.id AND p.member_id = :memberId AND p.deleted_at IS NULL) as isParticipant,
        (SELECT COUNT(*) FROM offline_conversation_participant p2 WHERE p2.conversation_id = c.id AND p2.deleted_at IS NULL) as numberOfParticipants
      FROM offline_conversation c
      WHERE c.id = :conversationId
      """, nativeQuery = true)
  Optional<OfflineConversationDetailProjection> findDetail(
      @Param("conversationId") UUID conversationId,
      @Param("memberId") UUID memberId
  );
}
