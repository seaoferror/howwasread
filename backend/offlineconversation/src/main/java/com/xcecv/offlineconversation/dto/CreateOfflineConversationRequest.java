package com.xcecv.offlineconversation.dto;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;

import java.time.Instant;

public record CreateOfflineConversationRequest(
    String novel,
    String poem,
    String shortStory,
    String play,
    String film,
    @NotBlank String by,
    @NotBlank String rule,
    @Min(value = 1) int capacity,
    @NotBlank Instant when,
    @NotBlank String where,
    String city,
    @Min(value = -90) @Max(value = 90) double lat,
    @Min(value = -180) @Max(value = 180) double lng,
    @NotBlank String h3Index
) {}