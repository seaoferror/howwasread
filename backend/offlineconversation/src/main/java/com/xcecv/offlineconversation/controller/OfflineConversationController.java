package com.xcecv.offlineconversation.controller;

import com.xcecv.offlineconversation.dto.CreateOfflineConversationRequest;
import com.xcecv.offlineconversation.dto.JoinOfflineConversationRequest;
import com.xcecv.offlineconversation.service.OfflineConversationService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/")
@RequiredArgsConstructor
public class OfflineConversationController {

  private final OfflineConversationService offlineConversationService;

  @PostMapping("/create")
  public ResponseEntity<?> create(
      @RequestBody CreateOfflineConversationRequest request,
      @RequestHeader("X-User-Id") UUID memberId) {
    Map<String, UUID> response = offlineConversationService.create(
        memberId,
        request.name(),
        request.lat(),
        request.lng()
    );
    return ResponseEntity.ok(response);
  }

  @PutMapping("/join")
  public ResponseEntity<?> join(
      @RequestBody JoinOfflineConversationRequest request,
      @RequestHeader("X-User-Id") UUID memberId
  ) {
    offlineConversationService.join(
        request.conversationId(),
        memberId
    );
    return ResponseEntity.ok("ok");
  }

  @GetMapping("/find/nearby")
  public ResponseEntity<?> findNearby(
      @RequestParam double latitude,
      @RequestParam double longitude,
      @RequestParam int resolution
  ) {
    if(resolution < 0 || resolution > 15){
      throw new ResponseStatusException(
          HttpStatus.BAD_REQUEST,
          "not acceptable resolution"
      );
    }
    offlineConversationService.findNearBy(latitude, longitude);

    return ResponseEntity.ok("ok");
  }
}
