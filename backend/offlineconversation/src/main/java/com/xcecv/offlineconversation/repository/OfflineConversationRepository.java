package com.xcecv.offlineconversation.repository;

import com.xcecv.offlineconversation.domain.OfflineConversation;
import com.xcecv.offlineconversation.projection.OfflineConversationMapProjection;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.time.Instant;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface OfflineConversationRepository extends JpaRepository<OfflineConversation, UUID> {
  List<OfflineConversationMapProjection> findTop2ByH3Res5AndTimeAfter(String h3Res5, Instant now);

  List<OfflineConversationMapProjection> findByH3Res7AndTimeAfter(String h3Res7, Instant now);

  <T> Optional<T> findById(UUID id, Class<T> type);
}
