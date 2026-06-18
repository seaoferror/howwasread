package com.xcecv.offlineconversation.dto;

import lombok.Builder;

import java.util.List;

@Builder
public record ChatMessage(
    byte[] id,
    byte[] fromId,
    String toIdType,
    byte[] toId,
    String contentType,
    List<String> contents
) {
}