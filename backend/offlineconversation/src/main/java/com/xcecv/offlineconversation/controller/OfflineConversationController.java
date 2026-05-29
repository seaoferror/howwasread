package com.xcecv.offlineconversation.controller;

import com.xcecv.offlineconversation.dto.CreateOfflineConversationRequest;
import com.xcecv.offlineconversation.dto.JoinOfflineConversationRequest;
import com.xcecv.offlineconversation.service.OfflineConversationService;
import jakarta.validation.Valid;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/")
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

  @PutMapping("/join")
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

  @GetMapping("/show")
  public ResponseEntity<?> showConvos(
      @NotBlank @RequestParam String h3Index,
      @NotNull @RequestHeader("X-User-Id") UUID memberId
  ) {
    var response = offlineConversationService.showConvos(h3Index, memberId);

    return ResponseEntity.ok(response);
  }
}
