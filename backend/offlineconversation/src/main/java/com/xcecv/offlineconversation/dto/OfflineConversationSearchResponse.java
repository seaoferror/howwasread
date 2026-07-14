package com.xcecv.offlineconversation.dto;

import lombok.Builder;
import org.springframework.data.elasticsearch.annotations.Field;
import org.springframework.data.elasticsearch.annotations.FieldType;

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
    String rule,
    String mapsLink,
    String location,
    Instant time,
    double lat,
    double lng
) {
}
