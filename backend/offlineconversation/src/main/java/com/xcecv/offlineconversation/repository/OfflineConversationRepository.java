package com.xcecv.offlineconversation.repository;

import com.xcecv.offlineconversation.domain.OfflineConversation;
import com.xcecv.offlineconversation.projection.FindByH3Result;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface OfflineConversationRepository extends JpaRepository<OfflineConversation, UUID> {
  List<FindByH3Result> findTop2ByH3Res5(String h3Res5);

  List<FindByH3Result> findByH3Res7(String h3Res7);
}
