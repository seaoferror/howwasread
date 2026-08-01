package com.xcecv.offlineconversation.dto;

import lombok.Builder;

import java.time.Instant;
import java.util.UUID;

@Builder
public record OfflineConversationDocument(
    UUID id,
    String novel,
    String poem,
    String shortStory,
    String play,
    String film,
    String writtenBy,
    Instant time,
    String h3Res5,
    String h3Res7,
    double latitude,
    double longitude
) {
}
