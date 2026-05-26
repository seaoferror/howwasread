package com.xcecv.offlineconversation.repository;

import com.xcecv.offlineconversation.domain.OfflineConversation;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface OfflineConversationRepository extends JpaRepository<OfflineConversation, UUID> {
  List<OfflineConversation> findByH3IndexIn(List<String> h3Indexes);
}
