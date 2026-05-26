package com.xcecv.offlineconversation.dto;

public record CreateOfflineConversationRequest(
    String name,
    double lat,
    double lng
) {}