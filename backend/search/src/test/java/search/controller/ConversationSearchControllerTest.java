package search.controller;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.http.HttpStatus;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.assertj.MockMvcTester;
import search.dto.OfflineConversationMapResponse;
import search.dto.OfflineConversationSearchResponse;
import search.dto.OnlineConversationSearchResponse;
import search.service.ConversationSearchService;

import java.time.Instant;
import java.util.List;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;

@WebMvcTest(ConversationSearchController.class)
class ConversationSearchControllerTest {

  private static final String BASE = "/search/conversation";

  @Autowired
  private MockMvcTester mvc;

  @MockitoBean
  private ConversationSearchService conversationSearchService;

  private final UUID id = UUID.randomUUID();
  private final Instant time = Instant.parse("2026-09-24T03:00:00Z");

  @Test
  void offline_parsesParamsAndReturnsJson() {
    when(conversationSearchService.searchOfflines("ham", "7", List.of("a", "b"), time, 1)).thenReturn(List.of(
        OfflineConversationSearchResponse.builder().id(id).novel("Hamlet").lat(37.5).lng(127.0).build()));

    assertThat(mvc.get().uri(BASE + "/offline")
        .param("input", "ham")
        .param("resolution", "7")
        .param("h3Indexes", "a", "b")
        .param("time", time.toString())
        .param("page", "1"))
        .hasStatusOk()
        .bodyJson().extractingPath("$[0].novel").isEqualTo("Hamlet");
  }

  @Test
  void offline_missingInput_isBadRequest() {
    assertThat(mvc.get().uri(BASE + "/offline")
        .param("resolution", "7")
        .param("h3Indexes", "a")
        .param("time", time.toString())
        .param("page", "1"))
        .hasStatus(HttpStatus.BAD_REQUEST);
    verifyNoInteractions(conversationSearchService);
  }

  @Test
  void online_returnsJson() {
    when(conversationSearchService.searchOnlines("ozy", time, 2)).thenReturn(List.of(
        OnlineConversationSearchResponse.builder().id(id).poem("Ozymandias").build()));

    assertThat(mvc.get().uri(BASE + "/online")
        .param("input", "ozy")
        .param("time", time.toString())
        .param("page", "2"))
        .hasStatusOk()
        .bodyJson().extractingPath("$[0].poem").isEqualTo("Ozymandias");
  }

  @Test
  void onlineList_clampsPageToOne() {
    when(conversationSearchService.listOnlines(time, 1)).thenReturn(List.of());

    assertThat(mvc.get().uri(BASE + "/online/list")
        .param("time", time.toString())
        .param("page", "0"))
        .hasStatusOk();
    verify(conversationSearchService).listOnlines(time, 1);
  }

  @Test
  void offlineMap_returnsJson() {
    when(conversationSearchService.mapOfflines("5", "h3", time)).thenReturn(List.of(
        OfflineConversationMapResponse.builder().id(id).writtenBy("shakespeare").lat(37.5).lng(127.0).build()));

    assertThat(mvc.get().uri(BASE + "/offline/map")
        .param("resolution", "5")
        .param("h3Index", "h3")
        .param("time", time.toString()))
        .hasStatusOk()
        .bodyJson().extractingPath("$[0].writtenBy").isEqualTo("shakespeare");
  }

  @Test
  void offlineMap_missingTime_isBadRequest() {
    assertThat(mvc.get().uri(BASE + "/offline/map")
        .param("resolution", "5")
        .param("h3Index", "h3"))
        .hasStatus(HttpStatus.BAD_REQUEST);
    verifyNoInteractions(conversationSearchService);
  }
}
