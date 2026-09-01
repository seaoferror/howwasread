package com.howwasread;

import com.howwasread.dto.Content;
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
    bufferedEvents.put(event.getKeyId(), event);
    ctx.timerService().registerProcessingTimeTimer(event.getScheduledTime());
  }

  @Override
  public void onTimer(long timestamp, KeyedProcessFunction<String, IncomingNotificationEvent, OutgoingNotificationEvent>.OnTimerContext ctx, Collector<OutgoingNotificationEvent> out) throws Exception {
    Map<UUID, Content> notifications = new HashMap<>();
    UUID partitionId = null;
    String partitionIdType = null;

    Iterator<Map.Entry<UUID, IncomingNotificationEvent>> iterator = bufferedEvents.entries().iterator();
    while (iterator.hasNext()) {
      IncomingNotificationEvent event = iterator.next().getValue();
      if (event.getScheduledTime() <= timestamp) {
        partitionId = event.getPartitionId();
        partitionIdType = event.getPartitionIdType();
        notifications.put(event.getKeyId(),
            Content.builder()
                .title(event.getTitle())
                .body(event.getBody())
                .build());
        iterator.remove();
      }
    }

    if (partitionId != null && !notifications.isEmpty()) {
      OutgoingNotificationEvent output = OutgoingNotificationEvent.builder()
          .partitionId(partitionId)
          .partitionIdType(partitionIdType)
          .notifications(notifications)
          .scheduledTime(timestamp)
          .build();
      out.collect(output);
    }
  }
}
