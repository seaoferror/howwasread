package com.xcecv.offlineconversation.service;

import com.uber.h3core.H3Core;
import com.xcecv.offlineconversation.domain.OfflineConversation;
import com.xcecv.offlineconversation.dto.FindNearByOfflineConversationResponse;
import com.xcecv.offlineconversation.repository.OfflineConversationRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.util.*;

@Service
@RequiredArgsConstructor
public class OfflineConversationService {

  private final H3Core h3;
  private final OfflineConversationRepository repository;

  public Map<String, UUID> create(
      UUID memberId,
      String name,
      double lat,
      double lng
  ) {
    String h3Index = h3.latLngToCellAddress(lat, lng, 15);
    var convo = OfflineConversation.builder()
        .name(name)
        .latitude(lat)
        .longitude(lng)
        .participants(new HashSet<>(Set.of(memberId)))
        .h3Index(h3Index)
        .build();

    return Map.of("id", repository.save(convo).getId());
  }

  public void join(UUID conversationId, UUID memberId) {
    var convo = repository.findById(conversationId)
        .orElseThrow(() -> new ResponseStatusException(
            HttpStatus.NOT_FOUND,
            "Conversation not found"
        ));

    convo.getParticipants().add(memberId);
    repository.save(convo);
  }

  public FindNearByOfflineConversationResponse findNearBy(
      double lat,
      double lng
  ) {

  }

  public List<OfflineConversation> findNearbyConversations(double userLat, double userLng, int resolution) {

    String originHex = h3.latLngToCellAddress(userLat, userLng, resolution);

    List<String> nearbyHexes = h3.gridDisk(originHex, 1);

    return repository.findByH3IndexIn(nearbyHexes);
  }


}
