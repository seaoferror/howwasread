package com.xcecv.offlineconversation.controller;

import com.xcecv.offlineconversation.service.OfflineConversationService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/offlineconversation")
@RequiredArgsConstructor
public class OfflineConversationController {
    private final OfflineConversationService offlineConversationService;

    
}
