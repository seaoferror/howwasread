package com.xcecv.offlineconversation.dto;

import java.util.UUID;

public record JoinOfflineConversationRequest(
    UUID conversationId
) {}
