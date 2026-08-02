package com.xcecv.search.dto;

import lombok.Builder;

import java.time.Instant;
import java.util.UUID;

@Builder
public record OnlineConversationSearchResponse(
    UUID id,
    String novel,
    String poem,
    String shortStory,
    String play,
    String film,
    String writtenBy,
    Instant time
) {
}
