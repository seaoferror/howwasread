package com.xcecv.offlineconversation.dto;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

import java.time.Duration;
import java.time.Instant;

public record CreateOfflineConversationRequest(
    String novel,
    String poem,
    String shortStory,
    String play,
    String film,
    @NotBlank String writtenBy,
    String rule,
    @NotNull Instant time,
    @Min(value = 0) int length,
    @NotBlank String googleMapsLink,
    @NotBlank String location,
    String city,
    @Min(value = -90) @Max(value = 90) double lat,
    @Min(value = -180) @Max(value = 180) double lng,
    @NotBlank String h3Res5,
    @NotBlank String h3Res7
) {
}