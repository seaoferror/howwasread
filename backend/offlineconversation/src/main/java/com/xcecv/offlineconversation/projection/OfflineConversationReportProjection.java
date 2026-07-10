package com.xcecv.offlineconversation.projection;

import java.util.Set;
import java.util.UUID;

public interface OfflineConversationReportProjection {
  Set<UUID> getReporterIds();
}
