package com.xcecv.offlineconversation.dto;


import lombok.Builder;

import java.util.UUID;

@Builder(toBuilder = true)
public record OfflineConversationMapResponse(
    UUID id,
    String writtenBy,
    double lat,
    double lng
) {
}
