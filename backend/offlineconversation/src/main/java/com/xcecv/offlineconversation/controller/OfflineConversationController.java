package com.xcecv.offlineconversation.controller;

import com.xcecv.offlineconversation.dto.CreateOfflineConversationRequest;
import com.xcecv.offlineconversation.dto.JoinOfflineConversationRequest;
import com.xcecv.offlineconversation.dto.OfflineConversationMapResponse;
import com.xcecv.offlineconversation.dto.OfflineConversationSearchResponse;
import com.xcecv.offlineconversation.service.OfflineConversationService;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/offlineconversation")
@RequiredArgsConstructor
public class OfflineConversationController {

  private final OfflineConversationService offlineConversationService;

  @PostMapping("/create")
  public ResponseEntity<?> create(
      @Valid @RequestBody CreateOfflineConversationRequest request,
      @NotNull @RequestHeader("X-User-Id") UUID memberId) {
    Map<String, UUID> response = offlineConversationService.create(request, memberId);
    return ResponseEntity.ok(response);
  }

  @PatchMapping("/join")
  public ResponseEntity<?> join(
      @Valid @RequestBody JoinOfflineConversationRequest request,
      @NotNull @RequestHeader("X-User-Id") UUID memberId
  ) {
    offlineConversationService.join(
        request,
        memberId
    );
    return ResponseEntity.ok("ok");
  }

  @PatchMapping("/quit")
  public ResponseEntity<?> quit(
      @Valid @RequestBody JoinOfflineConversationRequest request,
      @NotNull @RequestHeader("X-User-Id") UUID memberId
  ) {
    offlineConversationService.quit(
        request,
        memberId
    );
    return ResponseEntity.ok("ok");
  }

  @GetMapping("/map")
  public ResponseEntity<?> mapFarConvos(
      @NotBlank @RequestParam String resolution,
      @NotBlank @RequestParam String h3Index
  ) {
    List<OfflineConversationMapResponse> response = null;
    if (resolution.equals("5")) {
      response = offlineConversationService.mapRes5Convos(h3Index);
    }
    if (resolution.equals("7")) {
      response = offlineConversationService.mapRes7Convos(h3Index);
    }
    return ResponseEntity.ok(response);
  }

  @GetMapping("/detail")
  public ResponseEntity<?> detail(
      @NotNull @RequestParam UUID conversationId,
      @NotNull @RequestHeader("X-User-Id") UUID memberId
  ) {
    var response = offlineConversationService.detail(conversationId, memberId);
    return ResponseEntity.ok(response);
  }

  @PostMapping("/report")
  public ResponseEntity<?> report(
      @NotNull @RequestHeader("X-User-Id") UUID memberId,
      @Valid @RequestBody JoinOfflineConversationRequest request
  ) {
    offlineConversationService.report(request.conversationId(), memberId);
    return ResponseEntity.ok("ok");
  }

  @GetMapping("/search")
  public ResponseEntity<?> search(
      @NotBlank @RequestParam String input,
      @NotBlank @RequestParam String resolution,
      @NotBlank @RequestParam List<String> h3Indexes,
      @NotBlank @RequestParam int page
  ) {
    List<OfflineConversationSearchResponse> response = null;
    if (resolution.equals("5")) {
      response = offlineConversationService.searchH3Res5(input, h3Indexes, page);
    }
    if (resolution.equals("7")) {
      response = offlineConversationService.searchH3Res7(input, h3Indexes, page);
    }
    return ResponseEntity.ok(response);
  }
}
