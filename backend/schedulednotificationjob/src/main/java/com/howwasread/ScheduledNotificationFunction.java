package com.howwasread;

import com.howwasread.dto.IncomingNotificationEvent;
import com.howwasread.dto.OutgoingNotificationEvent;
import org.apache.flink.api.common.functions.OpenContext;
import org.apache.flink.api.common.state.MapState;
import org.apache.flink.api.common.state.MapStateDescriptor;
import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;

import java.io.Serial;
import java.util.*;

public class ScheduledNotificationFunction extends KeyedProcessFunction<String, IncomingNotificationEvent, OutgoingNotificationEvent> {

  @Serial
  private static final long serialVersionUID = 1L;

  private transient MapState<UUID, IncomingNotificationEvent> bufferedEvents;

  @Override
  public void open(OpenContext openContext) throws Exception {
    MapStateDescriptor<UUID, IncomingNotificationEvent> descriptor =
        new MapStateDescriptor<>("buffered-events", UUID.class, IncomingNotificationEvent.class);
    bufferedEvents = getRuntimeContext().getMapState(descriptor);
  }

  @Override
  public void processElement(IncomingNotificationEvent event, Context ctx, Collector<OutgoingNotificationEvent> out) throws Exception {
    bufferedEvents.put(event.getMemberId(), event);
    ctx.timerService().registerProcessingTimeTimer(event.getScheduledTime());
  }

  @Override
  public void onTimer(long timestamp, KeyedProcessFunction<String, IncomingNotificationEvent, OutgoingNotificationEvent>.OnTimerContext ctx, Collector<OutgoingNotificationEvent> out) throws Exception {
    List<UUID> memberIds = new ArrayList<>();
    String writtenBy = null;
    UUID conversationId = null;
    String notificationType = null;

    Iterator<Map.Entry<UUID, IncomingNotificationEvent>> iterator = bufferedEvents.entries().iterator();
    while (iterator.hasNext()) {
      IncomingNotificationEvent event = iterator.next().getValue();
      if (event.getScheduledTime() <= timestamp) {
        conversationId = event.getConversationId();
        writtenBy = event.getWrittenBy();
        notificationType = event.getNotificationType();
        memberIds.add(event.getMemberId());
        iterator.remove();
      }
    }

    if (conversationId != null && !memberIds.isEmpty()) {
      OutgoingNotificationEvent output = OutgoingNotificationEvent.builder()
          .conversationId(conversationId)
          .memberIds(memberIds)
          .writtenBy(writtenBy)
          .scheduledTime(timestamp)
          .notificationType(notificationType)
          .build();
      out.collect(output);
    }
  }
}
