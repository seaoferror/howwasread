package com.xcecv.offlineconversation.dto;


import lombok.Builder;

import java.time.Instant;
import java.util.UUID;

@Builder
public record ShowOfflineConversationResponse(
    UUID id,
    String novel,
    String poem,
    String shortStory,
    String play,
    String film,
    String writtenBy,
    String rule,
    int capacity,
    Instant time,
    String location,
    double lat,
    double lng,
    boolean isModerator,
    boolean isParticipant
) {
}
