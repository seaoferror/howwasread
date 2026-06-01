package com.xcecv.offlineconversation.projection;

import java.util.UUID;

public interface FindByH3Result {
  UUID getId();

  String getWrittenBy();

  double getLatitude();

  double getLongitude();
}
