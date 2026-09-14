package search.dto;

import lombok.Builder;

import java.time.Instant;
import java.util.UUID;

@Builder
public record OfflineConversationSearchResponse(
    UUID id,
    String novel,
    String poem,
    String shortStory,
    String play,
    String film,
    String writtenBy,
    Instant time,
    double lat,
    double lng
) {
}
