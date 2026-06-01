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

  @GetMapping("/map/far")
  public ResponseEntity<?> mapFarConvos(
      @NotBlank @RequestParam String h3Res5
  ) {
    var response = offlineConversationService.mapFarConvos(h3Res5);
    return ResponseEntity.ok(response);
  }

  @GetMapping("/map/close")
  public ResponseEntity<?> mapCloseConvos(
      @NotBlank @RequestParam String h3Res7
  ) {
    var response = offlineConversationService.mapCloseConvos(h3Res7);
    return ResponseEntity.ok(response);
  }
}
