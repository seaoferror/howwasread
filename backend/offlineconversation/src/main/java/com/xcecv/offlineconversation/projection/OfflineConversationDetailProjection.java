package com.xcecv.offlineconversation.projection;

import java.time.Duration;
import java.time.Instant;
import java.util.Set;
import java.util.UUID;

public interface OfflineConversationDetailProjection {
  String getNovel();

  String getPoem();

  String getShortStory();

  String getPlay();

  String getFilm();

  String getWrittenBy();

  String getRule();

  Instant getTime();

  Duration getLength();

  String getMapsLink();

  String getLocation();

  Set<UUID> getModeratorIds();

  Set<UUID> getParticipants();
}
