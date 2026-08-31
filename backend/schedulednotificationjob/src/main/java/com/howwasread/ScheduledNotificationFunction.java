package com.howwasread;

import com.howwasread.dto.IncomingNotificationEvent;
import com.howwasread.dto.OutgoingNotificationEvent;
import org.apache.flink.api.common.functions.OpenContext;
import org.apache.flink.api.common.state.ListState;
import org.apache.flink.api.common.state.ListStateDescriptor;
import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

public class ScheduledNotificationFunction extends KeyedProcessFunction<String, IncomingNotificationEvent, OutgoingNotificationEvent> {

  private static final long serialVersionUID = 1L;


  private transient ListState<IncomingNotificationEvent> bufferedEvents;

  @Override
  public void open(OpenContext openContext) throws Exception {
    ListStateDescriptor<IncomingNotificationEvent> descriptor = new ListStateDescriptor<>("buffered-events", IncomingNotificationEvent.class);
    bufferedEvents = getRuntimeContext().getListState(descriptor);
  }

  @Override
  public void processElement(IncomingNotificationEvent event, KeyedProcessFunction<String, IncomingNotificationEvent, OutgoingNotificationEvent>.Context ctx, Collector<OutgoingNotificationEvent> out) throws Exception {
    bufferedEvents.add(event);
    ctx.timerService().registerProcessingTimeTimer(event.getScheduledTime());
  }

  @Override
  public void onTimer(long timestamp, KeyedProcessFunction<String, IncomingNotificationEvent, OutgoingNotificationEvent>.OnTimerContext ctx, Collector<OutgoingNotificationEvent> out) throws Exception {

    List<UUID> memberIds = new ArrayList<>();
    List<IncomingNotificationEvent> remaining = new ArrayList<>();
    String writtenBy = null;
    UUID conversationId = null;
    String notificationType = null;
    for (IncomingNotificationEvent event : bufferedEvents.get()) {
      if (event.getScheduledTime() <= timestamp) {
        conversationId = event.getConversationId();
        writtenBy = event.getWrittenBy();
        notificationType = event.getNotificationType();
        if (event.getMemberId() != null) {
          memberIds.add(event.getMemberId());
        }
        continue;
      }
      remaining.add(event);
    }
    if (conversationId != null && !memberIds.isEmpty()) {
      OutgoingNotificationEvent output = OutgoingNotificationEvent.builder().conversationId(conversationId).memberIds(memberIds).writtenBy(writtenBy).scheduledTime(timestamp).notificationType(notificationType).build();
      out.collect(output);
    }
    if (remaining.isEmpty()) {
      bufferedEvents.clear();
      return;
    }
    bufferedEvents.update(remaining);
  }
}
