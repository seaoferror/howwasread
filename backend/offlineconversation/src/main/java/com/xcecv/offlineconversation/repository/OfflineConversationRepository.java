package com.xcecv.offlineconversation.repository;

import com.xcecv.offlineconversation.domain.OfflineConversation;
import com.xcecv.offlineconversation.projection.OfflineConversationMapProjection;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface OfflineConversationRepository extends JpaRepository<OfflineConversation, UUID> {
  List<OfflineConversationMapProjection> findTop2ByH3Res5(String h3Res5);

  List<OfflineConversationMapProjection> findByH3Res7(String h3Res7);

  <T> Optional<T> findById(UUID id, Class<T> type);
}
