package com.xcecv.offlineconversation.projection;

import lombok.Builder;

@Builder
public record OfflineConversationPinProjection(
    String writtenBy,
    double latitude,
    double longitude
) {
}