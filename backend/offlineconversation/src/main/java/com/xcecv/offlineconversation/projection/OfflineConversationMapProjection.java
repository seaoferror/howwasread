package com.xcecv.offlineconversation.projection;

import java.util.UUID;

public interface OfflineConversationMapProjection {
  UUID getId();

  String getWrittenBy();

  double getLatitude();

  double getLongitude();
}
