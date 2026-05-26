package com.xcecv.offlineconversation.service;

import com.uber.h3core.H3Core;
import com.xcecv.offlineconversation.domain.OfflineConversation;
import com.xcecv.offlineconversation.repository.OfflineConversationRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
public class OfflineConversationService {
  private final H3Core h3;
  private final OfflineConversationRepository repository;


  public OfflineConversation createConversation(double lat, double lng, int resolution) {
    OfflineConversation conversation = new OfflineConversation();
    conversation.setLatitude(lat);
    conversation.setLongitude(lng);

    String hexAddress = h3.latLngToCellAddress(lat, lng, resolution);
    conversation.setH3Index(hexAddress);

    return repository.save(conversation);
  }

  public List<OfflineConversation> findNearbyConversations(double userLat, double userLng, int resolution) {

    String originHex = h3.latLngToCellAddress(userLat, userLng, resolution);

    List<String> nearbyHexes = h3.gridDisk(originHex, 1);

    return repository.findByH3IndexIn(nearbyHexes);
  }
}
