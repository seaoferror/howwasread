package com.xcecv.search.controller;

import com.xcecv.search.dto.OfflineConversationSearchResponse;
import com.xcecv.search.service.SearchService;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.List;

@Controller
@RequiredArgsConstructor
@RequestMapping("/search")
public class SearchController {

  private final SearchService searchService;

  @GetMapping("/conversation/offline")
  public ResponseEntity<?> search(
      @NotBlank @RequestParam String input,
      @NotBlank @RequestParam String resolution,
      @NotEmpty @RequestParam List<String> h3Indexes,
      @RequestParam int page
  ) {
    List<OfflineConversationSearchResponse> response = null;
    if (resolution.equals("5")) {
      response = searchService.searchH3Res5(input, h3Indexes, page);
    }
    if (resolution.equals("7")) {
      response = searchService.searchH3Res7(input, h3Indexes, page);
    }
    return ResponseEntity.ok(response);
  }
}
