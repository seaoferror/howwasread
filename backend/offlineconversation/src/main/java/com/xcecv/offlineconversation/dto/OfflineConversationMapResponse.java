package com.xcecv.offlineconversation.dto;


import lombok.Builder;

import java.time.Instant;
import java.util.UUID;

@Builder
public record OfflineConversationMapResponse(
    UUID id,
    String writtenBy,
    double lat,
    double lng
) {
}
