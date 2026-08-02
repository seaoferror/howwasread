package com.xcecv.search.controller;

import com.xcecv.search.dto.OfflineConversationSearchResponse;
import com.xcecv.search.dto.OnlineConversationSearchResponse;
import com.xcecv.search.service.ConversationSearchService;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.time.Instant;
import java.util.List;

@Controller
@RequiredArgsConstructor
@RequestMapping("/search/conversation")
public class ConversationSearchController {

  private final ConversationSearchService conversationSearchService;

  @GetMapping("/offline")
  public ResponseEntity<?> searchOfflineConversations(
      @NotBlank @RequestParam String input,
      @NotBlank @RequestParam String resolution,
      @NotEmpty @RequestParam List<String> h3Indexes,
      @NotNull @RequestParam Instant time,
      @RequestParam int page
  ) {
    List<OfflineConversationSearchResponse> response = conversationSearchService.searchOfflines(input, resolution, h3Indexes, time, page);
    return ResponseEntity.ok(response);
  }

  @GetMapping("/online")
  public ResponseEntity<?> searchOnlineConversations(
      @NotBlank @RequestParam String input,
      @NotNull @RequestParam Instant time,
      @RequestParam int page
  ) {
    List<OnlineConversationSearchResponse> response = conversationSearchService.searchOnlines(input, time, page);
    return ResponseEntity.ok(response);
  }
}
