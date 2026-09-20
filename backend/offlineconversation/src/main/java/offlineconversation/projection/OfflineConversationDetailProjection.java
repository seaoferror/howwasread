package offlineconversation.projection;

import java.time.Instant;

public interface OfflineConversationDetailProjection {
  String getNovel();

  String getPoem();

  String getShortStory();

  String getPlay();

  String getFilm();

  String getWrittenBy();

  String getRule();

  Instant getTime();

  int getLengthMinutes();

  String getMapsLink();

  String getLocation();
}
