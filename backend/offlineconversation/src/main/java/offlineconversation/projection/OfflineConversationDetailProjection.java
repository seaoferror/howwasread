package offlineconversation.projection;

import java.time.Duration;
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

  Duration getLength();

  String getMapsLink();

  String getLocation();
}
