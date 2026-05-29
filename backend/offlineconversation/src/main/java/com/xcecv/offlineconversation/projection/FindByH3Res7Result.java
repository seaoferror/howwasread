package com.xcecv.offlineconversation.projection;

import java.time.Instant;
import java.util.Set;
import java.util.UUID;

public interface FindByH3Res7Result {
  UUID getId();
  String getNovel();
  String getPoem();
  String getShortStory();
  String getPlay();
  String getFilm();
  String getWrittenBy();
  String getRule();
  int getCapacity();
  Instant getTime();
  String getLocation();
  double getLatitude();
  double getLongitude();
  Set<UUID> getModeratorIds();
  Set<UUID> getParticipants();
}
