package com.xcecv.offlineconversation.dto;

import jakarta.validation.constraints.NotNull;

import java.util.UUID;

public record JoinOfflineConversationRequest(
    @NotNull UUID conversationId
) {
}
