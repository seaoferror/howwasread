package com.xcecv.offlineconversation.repository;

import com.xcecv.offlineconversation.domain.OfflineConversationDocument;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface OfflineConversationDocumentRepository extends CrudRepository<OfflineConversationDocument, UUID> {
}
