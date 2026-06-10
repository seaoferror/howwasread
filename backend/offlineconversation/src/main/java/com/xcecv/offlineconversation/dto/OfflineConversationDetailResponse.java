package com.xcecv.offlineconversation.dto;

import lombok.Builder;

import java.time.Instant;

@Builder
public record OfflineConversationDetailResponse(
    String novel,
    String poem,
    String shortStory,
    String play,
    String film,
    String writtenBy,
    String rule,
    Instant time,
    int length,
    String mapsLink,
    String location,
    boolean isModerator,
    boolean isParticipant,
    int numberOfParticipants
) {
}
