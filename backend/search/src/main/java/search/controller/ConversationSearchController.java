package search.controller;

import search.dto.OfflineConversationMapResponse;
import search.dto.OfflineConversationSearchResponse;
import search.dto.OnlineConversationSearchResponse;
import search.service.ConversationSearchService;
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

  @GetMapping("/online/list")
  public ResponseEntity<?> listOnlineConversations(
      @NotNull @RequestParam Instant time,
      @RequestParam int page
  ) {
    List<OnlineConversationSearchResponse> response = conversationSearchService.listOnlines(time, Math.max(page, 1));
    return ResponseEntity.ok(response);
  }

  @GetMapping("/offline/map")
  public ResponseEntity<?> mapOfflineConversations(
      @NotBlank @RequestParam String resolution,
      @NotBlank @RequestParam String h3Index,
      @NotNull @RequestParam Instant time
  ) {
    List<OfflineConversationMapResponse> response = conversationSearchService.mapOfflines(resolution, h3Index, time);
    return ResponseEntity.ok(response);
  }
}
